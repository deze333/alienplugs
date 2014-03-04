package mandrill

import (
)

//------------------------------------------------------------
// Model
//------------------------------------------------------------

type Mandrill struct {
    key string
}

//------------------------------------------------------------
// Constructor
//------------------------------------------------------------

// Creates new Mandrill connection.
func New(apikey string) (m *Mandrill, err error) {
    m = &Mandrill{key: apikey}
    if err = m.Ping(); err != nil {
        m = nil
    }
    return
}
