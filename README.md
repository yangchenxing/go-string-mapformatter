# go-string-mapformatter
String formatter using maps

##Usage

verb:

    %(NAME|GO_VERB)

GO_VERB is the verbs in goland.

##Example

    fmt.Println(mapformatter.MustFormat("Hello %(name|s), you owe me %(money|.2f) dollar.",
        map[string]interface{}{
            "name": "anyone",
            "money": 10.3,
        }))
    // Output: Hello anyone, you owe me 10.30 dollar.