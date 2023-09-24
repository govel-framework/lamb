package lamb

import (
	"github.com/govel-framework/lamb/evaluator"
	"github.com/govel-framework/lamb/internal"
	"github.com/govel-framework/lamb/object"

	"github.com/govel-framework/govel"
)

// Render renders a lamb template.
func Render(c *govel.Context, file string, vars map[string]interface{}) {
	if govel.Store != nil {
		// get all the cookies and check if the session is valid
		sessions := make(map[string]interface{})

		for _, cookie := range c.Request.Cookies() {
			session, err := govel.Store.Get(c.Request, cookie.Name)

			if err != nil {
				continue // it is not a valid session
			}

			sessions[cookie.Name] = session.Values
		}

		if vars == nil {
			vars = make(map[string]interface{})
		}

		vars["sessions"] = sessions
	}

	// load the file
	err := internal.LoadFile(file, vars, c.Buf, evaluator.Eval, *object.NewEnvironment())

	if err != nil {
		panic(err.Error())
	}

}
