package mandrill

import (
    "fmt"
)

//------------------------------------------------------------
// Users Calls
//------------------------------------------------------------

// Ping server to check the API key works.
func (m *Mandrill) Ping() (err error) {
    resp, err := post(
        MNDRL_USERS_PING,
        map[string]string{"key": m.key})

    if err != nil {
        return
    }

    // Must return 'PONG!'
    switch resp.(type) {
    case string:
        if resp.(string) == "PONG!" {
            return nil
        } else {
            return fmt.Errorf("Unknown ping response, assuming error")
        }
    default:
        return testError(resp)
    }
}

// Retrieve current user info.
func (m *Mandrill) UserInfo() (resp interface{}, err error) {
    resp, err = post(
        MNDRL_USERS_INFO,
        map[string]string{"key": m.key})

    if err != nil {
        return
    }
    return
}

