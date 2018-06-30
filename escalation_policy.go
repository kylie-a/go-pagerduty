package pagerduty

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

const (
	escPath = "/escalation_policies"
)

// EscalationRule is a rule for an escalation policy to trigger.
type EscalationRule struct {
	ID      string      `json:"id,omitempty"`
	Delay   uint        `json:"escalation_delay_in_minutes,omitempty"`
	Targets []APIObject `json:"targets"`
}

// EscalationPolicy is a collection of escalation rules.
type EscalationPolicy struct {
	APIObject
	Name            string           `json:"name,omitempty"`
	EscalationRules []EscalationRule `json:"escalation_rules,omitempty"`
	Services        []APIReference   `json:"services,omitempty"`
	NumLoops        uint             `json:"num_loops,omitempty"`
	Teams           []APIReference   `json:"teams,omitempty"`
	Description     string           `json:"description,omitempty"`
	RepeatEnabled   bool             `json:"repeat_enabled,omitempty"`
}

type EscalationPolicyResponse struct {
	APIResponse
}

func (r EscalationPolicyResponse) GetResource() (Resource, error) {
	var dest EscalationPolicy
	err := r.getResourceFromResponse(&dest)
	return dest, err
}

func NewEscalationPolicyResponse(resp *http.Response) EscalationPolicyResponse {
	return EscalationPolicyResponse{APIResponse{raw: resp, apiType: EscalationPolicyResourceType}}
}

// ListEscalationPoliciesResponse is the data structure returned from calling the ListEscalationPolicies API endpoint.
type ListEscalationPoliciesResponse struct {
	APIListObject
	EscalationPolicies []EscalationPolicy `json:"escalation_policies"`
}

type ListEscalationRulesResponse struct {
	APIListObject
	EscalationRules []EscalationRule `json:"escalation_rules"`
}

// ListEscalationPoliciesOptions is the data structure used when calling the ListEscalationPolicies API endpoint.
type ListEscalationPoliciesOptions struct {
	APIListObject
	Query    string   `url:"query,omitempty"`
	UserIDs  []string `url:"user_ids,omitempty,brackets"`
	TeamIDs  []string `url:"team_ids,omitempty,brackets"`
	Includes []string `url:"include,omitempty,brackets"`
	SortBy   string   `url:"sort_by,omitempty"`
}

// GetEscalationRuleOptions is the data structure used when calling the GetEscalationRule API endpoint.
type GetEscalationRuleOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// ListEscalationPolicies lists all of the existing escalation policies.
func (c *Client) ListEscalationPolicies(opts ...ResourceRequestOptionFunc) (*ListEscalationPoliciesResponse, error) {
	resp, err := c.ListResources(EscalationPolicyResourceType, opts...)
	if err != nil {
		return nil, err
	}
	var result ListEscalationPoliciesResponse
	return &result, deserialize(resp, &result)
}

// CreateEscalationPolicy creates a new escalation policy.
func (c *Client) CreateEscalationPolicy(e EscalationPolicy) (*EscalationPolicy, error) {
	resp, err := c.CreateResource(e)
	if err != nil {
		return nil, err
	}
	escPol := resp.(EscalationPolicy)
	return &escPol, nil
}

// DeleteEscalationPolicy deletes an existing escalation policy and rules.
func (c *Client) DeleteEscalationPolicy(id string) error {
	err := c.DeleteResource(EscalationPolicyResourceType, id)
	return err
}

// GetEscalationPolicyOptions is the data structure used when calling the GetEscalationPolicy API endpoint.
type GetEscalationPolicyOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// GetEscalationPolicy gets information about an existing escalation policy and its rules.
func (c *Client) GetEscalationPolicy(id string, opts ...ResourceRequestOptionFunc) (*EscalationPolicy, error) {
	res, err := c.GetResource(EscalationPolicyResourceType, id, opts...)
	if err != nil {
		return nil, err
	}
	escPol := res.(EscalationPolicy)
	return &escPol, nil
}

// UpdateEscalationPolicy updates an existing escalation policy and its rules.
func (c *Client) UpdateEscalationPolicy(id string, e *EscalationPolicy) (*EscalationPolicy, error) {
	data := make(map[string]EscalationPolicy)
	data["escalation_policy"] = *e
	resp, err := c.put(escPath+"/"+id, data, nil)
	return getEscalationPolicyFromResponse(c, resp, err)
}

// CreateEscalationRule creates a new escalation rule for an escalation policy
// and appends it to the end of the existing escalation rules.
func (c *Client) CreateEscalationRule(escID string, e EscalationRule) (*EscalationRule, error) {
	data := make(map[string]EscalationRule)
	data["escalation_rule"] = e
	resp, err := c.post(escPath+"/"+escID+"/escalation_rules", data)
	return getEscalationRuleFromResponse(c, resp, err)
}

// GetEscalationRule gets information about an existing escalation rule.
func (c *Client) GetEscalationRule(escID string, id string, o *GetEscalationRuleOptions) (*EscalationRule, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get(escPath + "/" + escID + "/escalation_rules/" + id + "?" + v.Encode())
	return getEscalationRuleFromResponse(c, resp, err)
}

// DeleteEscalationRule deletes an existing escalation rule.
func (c *Client) DeleteEscalationRule(escID string, id string) error {
	_, err := c.delete(escPath + "/" + escID + "/escalation_rules/" + id)
	return err
}

// UpdateEscalationRule updates an existing escalation rule.
func (c *Client) UpdateEscalationRule(escID string, id string, e *EscalationRule) (*EscalationRule, error) {
	data := make(map[string]EscalationRule)
	data["escalation_rule"] = *e
	resp, err := c.put(escPath+"/"+escID+"/escalation_rules/"+id, data, nil)
	return getEscalationRuleFromResponse(c, resp, err)
}

// ListEscalationRules lists all of the escalation rules for an existing escalation policy.
func (c *Client) ListEscalationRules(escID string) (*ListEscalationRulesResponse, error) {
	resp, err := c.get(escPath + "/" + escID + "/escalation_rules")
	if err != nil {
		return nil, err
	}

	var result ListEscalationRulesResponse
	return &result, deserialize(resp, &result)
}

func getEscalationRuleFromResponse(c *Client, resp *http.Response, err error) (*EscalationRule, error) {
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var target map[string]EscalationRule
	if dErr := deserialize(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	rootNode := "escalation_rule"
	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}
	return &t, nil
}

func getEscalationPolicyFromResponse(c *Client, resp *http.Response, err error) (*EscalationPolicy, error) {
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var target map[string]EscalationPolicy
	if dErr := deserialize(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	rootNode := "escalation_policy"
	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}
	return &t, nil
}
