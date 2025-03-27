package assets

import _ "embed"

//go:embed android_prompt.txt
var AndroidPrompt []byte

//go:embed ios_prompt.txt
var IOSPrompt []byte
