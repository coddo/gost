package email

import "fmt"

// ParseTemplate parses a standard HTML template and places the parameters in the
// template's indicated parts
func ParseTemplate(templateFormat string, params ...interface{}) string {
	return fmt.Sprintf(templateFormat, params)
}

const (
	activateAccountSubject  = `Welcome on board`
	activateAccountTemplate = `
<html>
<head>
    <style>
    p {
        font-size: 20px;
    }
    </style>
</head>
<body>
    <h2> Hello there,</h1> 
    <h2> Welcome to the GOST web framework!</h2>
    <br /> <br />
    <p>To activate your account, please use the following link: %s</p>
    <br /> <br />
    <p> Cheers! </p>
    <p> GOST Team </p>
</body>
</html>`
)

const (
	resetPasswordSubject  = `Password reset`
	resetPasswordTemplate = `
<html>
<head>
    <style>
    p {
        font-size: 20px;
    }
    </style>
</head>
<body>
    <h2> Hello there,</h1> <br/>
    <h2>A password reset was requested for this account. <br/><br/>
        If it was not you who made the request, please disregard this email
    </h2>
    <br /> <br />
    <p>To reset your password, please use the following link: %s</p>
    <br /> <br />
    <p> Cheers! </p>
    <p> GOST Team </p>
</body>
</html>`
)
