package main

import "flag"

type CMDFlags struct {
	isRewrite   bool
	isTranslate bool
	isSummarize bool
	isClipboard bool
	provider    string
	input       string
	language    string
	file        string
}

func SetFlags() *CMDFlags {
	flags := &CMDFlags{}

	var rewrite, r bool
	var translate, t bool
	var summarize, s bool
	var copyClipboard, c bool
	var provider, p string
	var input, i string
	var language, l string
	var file, f string

	flag.BoolVar(&rewrite, "rewrite", false, "AI rewrite function flag")
	flag.BoolVar(&r, "r", false, "AI rewrite function flag (shorthand)")

	flag.BoolVar(&translate, "translate", false, "AI translate function flag")
	flag.BoolVar(&t, "t", false, "AI translate function flag (shorthand)")

	flag.BoolVar(&summarize, "summarize", false, "AI summarize function flag")
	flag.BoolVar(&s, "s", false, "AI summarize function flag (shorthand)")

	flag.BoolVar(&copyClipboard, "clipboard", false, "Copy result to clipboard automatically")
	flag.BoolVar(&c, "c", false, "Copy result to clipboard automatically (shorthand)")

	flag.StringVar(&provider, "provider", "", "AI model provider flag")
	flag.StringVar(&p, "p", "", "AI model provider flag (shorthand)")

	flag.StringVar(&input, "input", "", "AI prompt")
	flag.StringVar(&i, "i", "", "AI prompt (shorthand)")

	flag.StringVar(&language, "language", "", "Translation target language")
	flag.StringVar(&l, "l", "", "Translation target language (shorthand)")

	flag.StringVar(&file, "file", "", "Use file as input")
	flag.StringVar(&f, "f", "", "Use file as input (shorthand)")

	flag.Parse()

	firstNonEmpty := func(a, b string) string {
		if a != "" {
			return a
		}
		return b
	}

	flags.isRewrite = rewrite || r
	flags.isTranslate = translate || t
	flags.isSummarize = summarize || s
	flags.isClipboard = copyClipboard || c
	flags.provider = firstNonEmpty(provider, p)
	flags.input = firstNonEmpty(input, i)
	flags.language = firstNonEmpty(language, l)
	flags.file = firstNonEmpty(file, f)

	return flags
}
