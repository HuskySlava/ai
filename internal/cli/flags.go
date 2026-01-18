package cli

import "flag"

type CMDFlags struct {
	IsRewrite   bool
	IsTranslate bool
	IsSummarize bool
	IsClipboard bool
	Provider    string
	Input       string
	Language    string
	File        string
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

	flags.IsRewrite = rewrite || r
	flags.IsTranslate = translate || t
	flags.IsSummarize = summarize || s
	flags.IsClipboard = copyClipboard || c
	flags.Provider = firstNonEmpty(provider, p)
	flags.Input = firstNonEmpty(input, i)
	flags.Language = firstNonEmpty(language, l)
	flags.File = firstNonEmpty(file, f)

	return flags
}
