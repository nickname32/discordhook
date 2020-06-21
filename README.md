# DiscordHook

[![Go.Dev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/nickname32/discordhook@v1.0.0?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/nickname32/discordhook)](https://goreportcard.com/report/github.com/nickname32/discordhook)

This module is made for those, who wants only use Discord webhooks. If you want more from Discord API, i would suggest you to use [disgord](https://github.com/andersfylling/disgord).

## Where to find webhook ID and token

You can find them in webhook URL
`https://discord.com/api/webhooks/WEBHOOK_ID_HERE/WEBHOOK_TOKEN_HERE`

## Examples

```Go
package main

import (
    "fmt"

    "github.com/nickname32/discordhook"
)

func main() {
    wa, err := discordhook.NewWebhookAPI(12345678900987654321, "TOKENtoken1234567890asdfghjkl", true, nil)
    if err != nil {
        panic(err)
    }

    wh, err := wa.Get(nil)
    if err != nil {
        panic(err)
    }

    fmt.Println(wh.Name)

    msg, err := wa.Execute(nil, &discordhook.WebhookExecuteParams{
        Content: "Example text",
        Embeds: []*discordhook.Embed{
            {
                Title:       "Hi",
                Description: "This is description",
            },
        },
    }, nil, "")
    if err != nil {
        panic(err)
    }

    fmt.Println(msg.ID)

    wh, err = wa.Modify(nil, &discordhook.WebhookModifyParams{
        Name: "This is new default webhook name",
    })
    if err != nil {
        panic(err)
    }

    fmt.Println(wh)

    err = wa.Delete(nil)
    if err != nil {
        panic(err)
    }
}
```
