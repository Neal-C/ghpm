// package for everything related to Github Privacy Management
package ghpm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"log"
	"net/http"
	"slices"
	"sync"
)

// STARS_THRESHOLD : the required numbers of stars on a repository for it be avoided by ghpm
const STARS_THRESHOLD uint = 1

type GithubPrivacyManager struct {
	// token that allows requesting github on behalf of the user
	githubAuthToken string
	// httpClient that does the requests
	httpClient *http.Client
	// the username for the user that did the oauth authentication process
	username string
}

type User struct {
	Username string `json:"login"`
}

type GithubRepository struct {
	Stars uint `json:"stargazers_count"`

	Fullname string `json:"full_name"`

	Private bool `json:"private"`

	IsFork bool `json:"fork"`
}

func Prettyfy(data any) (string, error) {
	val, err := json.MarshalIndent(data, "", "")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func ToFullname(repositories []GithubRepository) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, repo := range repositories {
			yield(repo.Fullname)
		}
	}
}

func NewGithubPrivacyManager(githubAuthToken string, httpClient *http.Client) GithubPrivacyManager {

	httpRequest, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.github.com/user", http.NoBody)

	// Authorization
	httpRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %s", githubAuthToken))

	// Recommended in the github API documentation
	httpRequest.Header.Add("Accept", "application/vnd.github+json")

	// Targeted Github API version
	httpRequest.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	httpResponse, err := httpClient.Do(httpRequest)

	if err != nil {
		log.Fatal("could not fetch username of for the given githubAuthToken. Please complain to the developer")
	}

	defer httpResponse.Body.Close()

	var user User

	if err := json.NewDecoder(httpResponse.Body).Decode(&user); err != nil {
		log.Fatal("could not get the login name (aka username) of for the given githubAuthToken. Please complain to the developer")
	}

	return GithubPrivacyManager{
		githubAuthToken: githubAuthToken,
		httpClient:      httpClient,
		username:        user.Username,
	}
}

func (self *GithubPrivacyManager) setRequiredHeadersOnGithubRequest(httpRequest *http.Request) {

	// Authorization
	httpRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", self.githubAuthToken))

	// Recommended in the github API documentation
	httpRequest.Header.Set("Accept", "application/vnd.github+json")

	// Recommended in the github API documentation
	httpRequest.Header.Set("User-Agent", "ghpm")

	// Targeted Github API version
	httpRequest.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}

func (self *GithubPrivacyManager) ListAllPublicRepositories(ctx context.Context) error {
	githubAPIEndpoint := "https://api.github.com/user/repos?visibility=public&per_page=100"

	httpRequest, _ := http.NewRequestWithContext(ctx, http.MethodGet, githubAPIEndpoint, http.NoBody)

	self.setRequiredHeadersOnGithubRequest(httpRequest)

	httpResponse, err := self.httpClient.Do(httpRequest)

	if err != nil {
		return err
	}

	defer httpResponse.Body.Close()

	switch {
	case httpResponse.StatusCode == http.StatusNotFound:

		return fmt.Errorf("%d : not found. Please complain to the developer", http.StatusNotFound)

	case httpResponse.StatusCode >= 500:

		return fmt.Errorf("github is likely down. Retry. If it does persist: Please complain to the developer")

	}

	var publicRepositories []GithubRepository

	if err := json.NewDecoder(httpResponse.Body).Decode(&publicRepositories); err != nil {
		return err
	}

	// var namesOfPublicRepositories = make([]string, 0, 100)

	// for _, repo := range publicRepositories {
	// 	namesOfPublicRepositories = append(namesOfPublicRepositories, repo.Fullname)
	// }

	namesOfPublicRepositories := slices.Collect(ToFullname(publicRepositories))

	names, err := Prettyfy(namesOfPublicRepositories)

	if err != nil {
		return err
	}

	fmt.Printf("your public repositories : %s \n", names)

	return nil
}

func (self *GithubPrivacyManager) ListAllPrivateRepositories(ctx context.Context) error {

	githubAPIEndpoint := "https://api.github.com/user/repos?visibility=private&per_page=100"

	httpRequest, _ := http.NewRequestWithContext(ctx, http.MethodGet, githubAPIEndpoint, http.NoBody)

	self.setRequiredHeadersOnGithubRequest(httpRequest)

	httpResponse, err := self.httpClient.Do(httpRequest)

	if err != nil {
		return err
	}

	defer httpResponse.Body.Close()

	switch {
	case httpResponse.StatusCode == http.StatusNotFound:

		return fmt.Errorf("%d : not found. Please complain to the developer", http.StatusNotFound)

	case httpResponse.StatusCode >= 500:

		return fmt.Errorf("github is likely down. Retry. If it does persist: Please complain to the developer")
	}

	var privateRepositories []GithubRepository

	if err := json.NewDecoder(httpResponse.Body).Decode(&privateRepositories); err != nil {
		return err
	}

	// var namesOfPrivateRepositories = make([]string, 0, 100)

	// for _, repo := range privateRepositories {
	// 	namesOfPrivateRepositories = append(namesOfPrivateRepositories, repo.Fullname)
	// }

	namesOfPrivateRepositories := slices.Collect(ToFullname(privateRepositories))

	names, err := Prettyfy(namesOfPrivateRepositories)

	if err != nil {
		return err
	}

	fmt.Printf("your private repositories : %s \n", names)

	return nil

}

func (self *GithubPrivacyManager) SwitchRepoToPrivateByName(ctx context.Context, repositoryName string) error {

	readmeRepository := fmt.Sprintf("%s/%s", self.username, self.username)

	targetRepository := fmt.Sprintf("%s/%s", self.username, repositoryName)

	if targetRepository == readmeRepository {
		return fmt.Errorf("it makes no sense to make private your %s.\nGo through the web ui for that", readmeRepository)
	}

	publicRepoEndpoint := fmt.Sprintf("https://api.github.com/repos/%s", targetRepository)

	httpRequest, _ := http.NewRequestWithContext(ctx, http.MethodGet, publicRepoEndpoint, http.NoBody)

	self.setRequiredHeadersOnGithubRequest(httpRequest)

	httpReponse, err := self.httpClient.Do(httpRequest)

	if err != nil {
		return err
	}

	var publicRepository GithubRepository

	if err := json.NewDecoder(httpReponse.Body).Decode(&publicRepository); err != nil {
		return err
	}

	if publicRepository.Stars >= STARS_THRESHOLD {
		return fmt.Errorf("repository cannot be switched to private by ghpm because it has more than %d ", STARS_THRESHOLD)
	}

	payload := map[string]any{
		"private": true,
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	httpPatchRequest, _ := http.NewRequestWithContext(ctx, http.MethodPatch, publicRepoEndpoint, bytes.NewBuffer(jsonPayload))

	self.setRequiredHeadersOnGithubRequest(httpPatchRequest)

	httpResponse, err := self.httpClient.Do(httpPatchRequest)

	if err != nil {
		return err
	}

	defer httpResponse.Body.Close()

	switch {
	case httpResponse.StatusCode == http.StatusUnprocessableEntity:

		return fmt.Errorf("repository %s was not switched to private. Consider using the web ui for this one", repositoryName)

	case httpResponse.StatusCode == http.StatusNotFound:

		return fmt.Errorf("repository %s was not switched to private because it was not found. Did you misspell?", repositoryName)
	}

	return nil

}

func (self *GithubPrivacyManager) SwitchRepoToPublicByName(ctx context.Context, repositoryName string) error {

	readmeRepository := fmt.Sprintf("%s/%s", self.username, self.username)

	targetRepository := fmt.Sprintf("%s/%s", self.username, repositoryName)

	if targetRepository == readmeRepository {

		return fmt.Errorf("it makes no sense to change your %s. It's your profile's README: it's meant to be read.\nGo through the web ui for that", readmeRepository)

	}

	privateRepositoryEndpoint := fmt.Sprintf("https://api.github.com/repos/%s", targetRepository)

	payload := map[string]any{
		"private": false,
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPatch, privateRepositoryEndpoint, bytes.NewBuffer(jsonPayload))

	if err != nil {
		return err
	}

	self.setRequiredHeadersOnGithubRequest(httpRequest)

	httpResponse, err := self.httpClient.Do(httpRequest)

	if err != nil {

		log.Printf("%s was not switched to private. I suggest to you try from the web version for this one. I am sorry for failing you, please complain to the developer", targetRepository)

		return err
	}

	defer httpResponse.Body.Close()

	switch {
	case httpResponse.StatusCode == http.StatusUnprocessableEntity:

		return fmt.Errorf("repository %s was not switched to public. Consider using the web ui for this one", repositoryName)

	case httpResponse.StatusCode == http.StatusNotFound:

		return fmt.Errorf("repository %s was not switched to public because it was not found. Did you misspell?", repositoryName)

	case httpResponse.StatusCode >= 500:

		return fmt.Errorf("github is likely down. Retry. If it does persist: Please complain to the developer")
	}

	return nil

}

func (self *GithubPrivacyManager) SwitchAllRepositoriesToPrivate(ctx context.Context) error {

	publicRepositoriesGithubAPIEndpoint := fmt.Sprintf("https://api.github.com/users/%s/repos?visibility=public&per_page=100", self.username)

	readmeRepository := fmt.Sprintf("%s/%s", self.username, self.username)

	for {

		publicRepositoriesHTTPRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, publicRepositoriesGithubAPIEndpoint, http.NoBody)

		if err != nil {
			return err
		}

		self.setRequiredHeadersOnGithubRequest(publicRepositoriesHTTPRequest)

		httpResponse, err := self.httpClient.Do(publicRepositoriesHTTPRequest)

		if err != nil {
			return err
		}

		switch {
		case httpResponse.StatusCode == http.StatusNotFound:

			return fmt.Errorf("%d : not found. Did you spell that right? that its name?", http.StatusNotFound)

		case httpResponse.StatusCode >= 500:

			return fmt.Errorf("github is likely down. Retry. If it does persist: Please complain to the developer")
		}

		var publicRepositories []GithubRepository

		if err := json.NewDecoder(httpResponse.Body).Decode(&publicRepositories); err != nil {
			return err
		}

		httpResponse.Body.Close()

		payload := map[string]any{
			"private": true,
		}

		jsonPayload, err := json.Marshal(payload)

		if err != nil {
			return fmt.Errorf("json.Marshal: %s", err)
		}

		var switchWaitGroup sync.WaitGroup

		// TODO : lobby github for a batch request endpoint, so that it can be only 1 HTTP call and not O(n) HTTP calls
		for _, repo := range publicRepositories {

			if repo.Fullname == readmeRepository {

				fmt.Printf("skipped %s because it's a special repository \n", readmeRepository)

				continue
			}

			if repo.Stars >= STARS_THRESHOLD {

				log.Printf("repository %s cannot be switched to private by ghpm because it has more than %d stars -> (%d) \n", repo.Fullname, STARS_THRESHOLD, repo.Stars)

				continue
			}

			if repo.IsFork {

				log.Printf("skipped %s because it's a fork \n", repo.Fullname)

				continue
			}

			switchWaitGroup.Add(1)

			go func() {

				defer switchWaitGroup.Done()

				currentPublicRepositoryEndpoint := fmt.Sprintf("https://api.github.com/repos/%s", repo.Fullname)

				httpPatchRequest, err := http.NewRequestWithContext(ctx, http.MethodPatch, currentPublicRepositoryEndpoint, bytes.NewBuffer(jsonPayload))

				if err != nil {

					log.Printf("error requesting %s: %s \n", repo.Fullname, err)
					log.Println("skipping", repo.Fullname)

					switchWaitGroup.Done()

					return
				}

				self.setRequiredHeadersOnGithubRequest(httpPatchRequest)

				httpResponse, err := self.httpClient.Do(httpPatchRequest)

				if err != nil {

					log.Printf("error processing %s; err=%s", repo.Fullname, err)

					switchWaitGroup.Done()

					return
				}

				httpResponse.Body.Close()

				switch {
				case httpResponse.StatusCode == http.StatusNotImplemented:

					log.Printf("%s was not switched to private. I suggest to you try from the web version for this one. I am sorry for failing you, please complain to the developer \n", repo.Fullname)

				case httpResponse.StatusCode == http.StatusNotFound:

					log.Printf("%s was not found. Did you spell that right? that its name? \n", repo.Fullname)

				case httpResponse.StatusCode >= 500:

					log.Printf("github is likely down. Retry. If it does persist: Please complain to the developer. %s not switched \n", repo.Fullname)

				default:

					log.Printf("%s switched to private. \n", repo.Fullname)
				}

			}()

		}

		switchWaitGroup.Wait()

		if len(publicRepositories) != 100 {
			break
		}

	}

	return nil

}
