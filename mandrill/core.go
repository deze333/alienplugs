package mandrill

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "bytes"
)

//------------------------------------------------------------
// Core functions
//------------------------------------------------------------

// Sends request to Mandirill.
func post(cmd string, params interface{}) (resp interface{}, err error) {

    // Build JSON request
    pjson, err := json.Marshal(params)
    if err != nil {
        return
    }

    // Send request
    rs, err := http.Post(
        MNDRL_API_URL + cmd,
        "application/json",
        bytes.NewBuffer(pjson))

    if err != nil {
        return
    }

    defer rs.Body.Close()

    // Response body
    body, err := ioutil.ReadAll(rs.Body)
    if err != nil {
        err = fmt.Errorf("Error reading Mandrill response body: %v", err)
        return
    }

    // DEBUG
    fmt.Println(string(body))

    // Unmarshal response into map
    err = json.Unmarshal(body, &resp)

    // Check response code
    if rs.StatusCode != 200 {
        err = testError(resp)
        return
    }

    // Test for possible error
    err = testError(resp)
    return
}

// Tests if response has error structure.
func testError(resp interface{}) (err error) {
    switch resp.(type) {
    case string:
        return
    case int:
        return
    case []interface{}:
        return
    case map[string]interface{}:
        m := resp.(map[string]interface{})
        if m["status"] == "error" {
            return fmt.Errorf("Mandrill Error: %v, %v, %v", 
                m["code"],
                m["name"],
                m["message"])
        } else {
            return
        }
    default:
        return
    }
}
