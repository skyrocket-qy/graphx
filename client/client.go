package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/skyrocketOoO/zanazibar-dag/domain"
)

type ZanzibarDagClient struct {
	Url string
}

func NewZanzibarDagClient(url string) (*ZanzibarDagClient, error) {
	resp, err := http.Get(fmt.Sprintf("%s/healthy", url))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("can not connect to zanzibar-dag application")
	}

	return &ZanzibarDagClient{
		Url: fmt.Sprintf("%s/relations", url),
	}, nil
}

func (r *ZanzibarDagClient) GetAll() ([]domain.Relation, error) {
	resp, err := http.Get(r.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	type respBody struct {
		Data []domain.Relation `json:"data"`
	}
	body := respBody{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Data, nil
}

func (r *ZanzibarDagClient) Query(relation domain.Relation) ([]domain.Relation, error) {
	resp, err := http.Get(r.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	type respBody struct {
		Data []domain.Relation `json:"data"`
	}
	body := respBody{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Data, nil
}

func (r *ZanzibarDagClient) Create(relation domain.Relation, existOk bool) error {
	type requestBody struct {
		Relation domain.Relation `json:"relation"`
		ExistOk  bool            `json:"exist_ok"`
	}
	reqBody := requestBody{
		Relation: relation,
		ExistOk:  existOk,
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(r.Url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (r *ZanzibarDagClient) Delete(relation domain.Relation) error {
	payload, err := json.Marshal(relation)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", r.Url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (r *ZanzibarDagClient) DeleteByQueries(queries []domain.Relation) error {
	type requestBody struct {
		Queries []domain.Relation `json:"queries"`
	}
	payload, err := json.Marshal(requestBody{Queries: queries})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", r.Url+"/delete-by-queries", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (r *ZanzibarDagClient) BatchOperation(operations []domain.Operation) error {
	type requestBody struct {
		Operations []domain.Operation `json:"operations"`
	}
	payload, err := json.Marshal(requestBody{Operations: operations})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", r.Url+"/batch-operation", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (r *ZanzibarDagClient) GetAllNamespaces() ([]string, error) {
	req, err := http.NewRequest("GET", r.Url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the GET request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response body to extract namespaces
	var namespacesResponse struct {
		Data []string `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&namespacesResponse); err != nil {
		return nil, err
	}

	return namespacesResponse.Data, nil
}

func (r *ZanzibarDagClient) Check(from domain.Node, to domain.Node, searchCond domain.SearchCondition) (bool, error) {
	type requestBody struct {
		Subject         domain.Node            `json:"subject"`
		Object          domain.Node            `json:"object"`
		SearchCondition domain.SearchCondition `json:"search_condition"`
	}
	payload := requestBody{
		Subject:         from,
		Object:          to,
		SearchCondition: searchCond,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", r.Url+"/check", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return true, nil
}

func (r *ZanzibarDagClient) GetShortestPath(from domain.Node, to domain.Node, searchCond domain.SearchCondition) ([]domain.Relation, error) {
	type requestBody struct {
		Subject         domain.Node            `json:"subject"`
		Object          domain.Node            `json:"object"`
		SearchCondition domain.SearchCondition `json:"search_condition"`
	}
	payload := requestBody{
		Subject:         from,
		Object:          to,
		SearchCondition: searchCond,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", r.Url+"/get-shortest-path", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	type respBody struct {
		Data []domain.Relation `json:"data"`
	}
	body := respBody{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Data, nil
}

func (r *ZanzibarDagClient) GetAllPaths(from domain.Node, to domain.Node, searchCond domain.SearchCondition) ([][]domain.Relation, error) {
	type requestBody struct {
		Subject         domain.Node            `json:"subject"`
		Object          domain.Node            `json:"object"`
		SearchCondition domain.SearchCondition `json:"search_condition"`
	}
	payload := requestBody{
		Subject:         from,
		Object:          to,
		SearchCondition: searchCond,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", r.Url+"/get-all-paths", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	type respBody struct {
		Data [][]domain.Relation `json:"data"`
	}
	body := respBody{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Data, nil
}

func (r *ZanzibarDagClient) GetAllObjectRelations(subject domain.Node, searchCond domain.SearchCondition, collectCond domain.CollectCondition) ([]domain.Relation, error) {
	type requestBody struct {
		Subject          domain.Node             `json:"subject"`
		SearchCondition  domain.SearchCondition  `json:"search_condition"`
		CollectCondition domain.CollectCondition `json:"collect_condition"`
	}
	payload := requestBody{
		Subject:          subject,
		SearchCondition:  searchCond,
		CollectCondition: collectCond,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", r.Url+"/get-all-object-relations", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	type respBody struct {
		Data []domain.Relation `json:"data"`
	}
	body := respBody{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Data, nil
}

func (r *ZanzibarDagClient) GetAllSubjectRelations(object domain.Node, searchCond domain.SearchCondition, collectCond domain.CollectCondition) ([]domain.Relation, error) {
	type requestBody struct {
		Object           domain.Node             `json:"object"`
		SearchCondition  domain.SearchCondition  `json:"search_condition"`
		CollectCondition domain.CollectCondition `json:"collect_condition"`
	}
	payload := requestBody{
		Object:           object,
		SearchCondition:  searchCond,
		CollectCondition: collectCond,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", r.Url+"/get-all-subject-relations", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	type respBody struct {
		Data []domain.Relation `json:"data"`
	}
	body := respBody{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Data, nil
}

func (r *ZanzibarDagClient) ClearAllRelations() error {

	req, err := http.NewRequest("POST", r.Url+"/clear-all-relations", nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
