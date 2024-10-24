package ghpm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
}

func Prettyfy(data any) (string, error) {
	val, err := json.MarshalIndent(data, "", "")
	if err != nil {
		return "", err
	}
	return string(val), nil
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

	if httpResponse.StatusCode == http.StatusNotFound {
		return fmt.Errorf("%d : not found. Please complain to the developer", http.StatusNotFound)
	}

	if httpResponse.StatusCode >= 500 {
		return fmt.Errorf("github is likely down. Retry. If it does persist: Please complain to the developer")
	}

	var publicRepositories []GithubRepository

	if err := json.NewDecoder(httpResponse.Body).Decode(&publicRepositories); err != nil {
		return err
	}

	var namesOfPublicRepositories []string

	for _, repo := range publicRepositories {
		namesOfPublicRepositories= append(namesOfPublicRepositories, repo.Fullname)
	}

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

	if httpResponse.StatusCode == http.StatusNotFound {
		return fmt.Errorf("%d : not found. Please complain to the developer", http.StatusNotFound)
	}

	if httpResponse.StatusCode >= 500 {
		return fmt.Errorf("github is likely down. Retry. If it does persist: Please complain to the developer")
	}

	var privateRepositories []GithubRepository

	if err := json.NewDecoder(httpResponse.Body).Decode(&privateRepositories); err != nil {
		return err
	}

	var namesOfPrivateRepositories []string

	for _, repo := range privateRepositories {
		namesOfPrivateRepositories = append(namesOfPrivateRepositories, repo.Fullname)
	}

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
		return fmt.Errorf("it makes no sense to make private your %s. \n Go through the web ui for that", readmeRepository)
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

	if httpResponse.StatusCode == http.StatusUnprocessableEntity {
		return fmt.Errorf("repository %s was not switched to private. Consider using the web ui for this one", repositoryName)
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		return fmt.Errorf("repository %s was not switched to private because it was not found", repositoryName)
	}

	return nil

}

func (self *GithubPrivacyManager) SwitchRepoToPublicByName(ctx context.Context, repositoryName string) error {

	readmeRepository := fmt.Sprintf("%s/%s", self.username, self.username)

	targetRepository := fmt.Sprintf("%s/%s", self.username, repositoryName)

	if targetRepository == readmeRepository {
		return fmt.Errorf("it makes no sense to make private your %s. It's your profile's README: it's meant to be read.\nGo through the web ui for that", readmeRepository)
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

	if httpResponse.StatusCode == http.StatusUnprocessableEntity {
		return fmt.Errorf("repository %s was not switched to public. Consider using the web ui for this one", repositoryName)
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		return fmt.Errorf("repository %s was not switched to public because it was not found", repositoryName)
	}

	if httpResponse.StatusCode >= 500 {
		return fmt.Errorf("github is likely down. Retry. If it does persist: Please complain to the developer")
	}

	defer httpResponse.Body.Close()

	return nil

}

func (self *GithubPrivacyManager) SwitchAllRepositoriesToPrivate(ctx context.Context) error {

	var shouldRunAgain bool

run:

	publicRepositoriesGithubAPIEndpoint := fmt.Sprintf("https://api.github.com/users/%s/repos?visibility=public&per_page=100", self.username)

	publicRepositoriesHTTPRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, publicRepositoriesGithubAPIEndpoint, http.NoBody)

	if err != nil {
		return err
	}

	self.setRequiredHeadersOnGithubRequest(publicRepositoriesHTTPRequest)

	httpResponse, err := self.httpClient.Do(publicRepositoriesHTTPRequest)

	if err != nil {
		return err
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		return fmt.Errorf("%d : not found. Please complain to the developer", http.StatusNotFound)
	}

	if httpResponse.StatusCode >= 500 {
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

	// TODO : lobby github for a batch request endpoint, so that it can be only 1 HTTP call and not O(n) HTTP calls
	for _, repo := range publicRepositories {

		if repo.Stars >= STARS_THRESHOLD {

			log.Printf("repository %s cannot be switched to private by ghpm because it has more than %d stars -> (%d) \n", repo.Fullname, STARS_THRESHOLD, repo.Stars)

			continue
		}

		readmeRepository := fmt.Sprintf("%s/%s", self.username, self.username)

		if repo.Fullname == readmeRepository {
			continue
		}

		currentPublicRepositoryEndpoint := fmt.Sprintf("https://api.github.com/repos/%s", repo.Fullname)

		jsonPayload, err := json.Marshal(payload)

		if err != nil {

			log.Printf("error processing %s; err=%s", repo.Fullname, err)

			continue
		}

		httpPatchRequest, _ := http.NewRequestWithContext(ctx, http.MethodPatch, currentPublicRepositoryEndpoint, bytes.NewBuffer(jsonPayload))

		self.setRequiredHeadersOnGithubRequest(httpPatchRequest)

		httpResponse, err := self.httpClient.Do(httpPatchRequest)

		if err != nil {

			log.Printf("error processing %s; err=%s", repo.Fullname, err)

			continue
		}

		switch httpResponse.StatusCode {
		case http.StatusNotImplemented, http.StatusNotFound:

			log.Printf("%s was not switched to private. I suggest to you try from the web version for this one. I am sorry for failing you, please complain to the developer \n", repo.Fullname)

			httpResponse.Body.Close()

			continue
		}

		if httpResponse.StatusCode >= 500 {

			log.Printf("github is likely down. Retry. If it does persist: Please complain to the developer \n")

			httpResponse.Body.Close()

			continue
		}

		log.Printf("%s switched to private. \n", repo.Fullname)

		httpResponse.Body.Close()

	}

	shouldRunAgain = len(publicRepositories) == 100

	if shouldRunAgain {
		goto run
	}

	return nil

}
