package css

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// parseBackgroundSize reads a comma-separated background-size list: each value
// is "cover", "contain" or one or two lengths, the missing axis becoming auto.
func parseBackgroundSize(raw string, ctx Units) ([]BackSize, bool) {
	var out []BackSize
	for _, part := range splitFields(raw, ',') {
		sz, ok := parseBackSizeOne(strings.TrimSpace(part), ctx)
		if !ok {
			return nil, false
		}
		out = append(out, sz)
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}

// parseBackSizeOne reads one background-size value.
func parseBackSizeOne(part string, ctx Units) (BackSize, bool) {
	tokens := splitTokens(part)
	if len(tokens) == 0 {
		return BackSize{}, false
	}
	if len(tokens) == 1 {
		switch tokens[0] {
		case "cover":
			return BackSize{Cover: true}, true
		case "contain":
			return BackSize{Contain: true}, true
		}
	}
	if len(tokens) > 2 {
		return BackSize{}, false
	}
	var sz BackSize
	var vals []Length
	for _, tk := range tokens {
		l, err := parseLengthAt(tk, ctx)
		if err != nil {
			return BackSize{}, false
		}
		vals = append(vals, l)
	}
	sz.W = vals[0]
	if len(vals) > 1 {
		sz.H = vals[1]
	} else {
		sz.H = Auto()
	}
	return sz, true
}

// parseBackgroundRepeat reads a comma-separated background-repeat list: one
// keyword for both axes, or an X and a Y keyword, with the shorthand
// repeat-x and repeat-y folded into their pair.
func parseBackgroundRepeat(raw string) ([]BackRepeat, bool) {
	var out []BackRepeat
	for _, part := range splitFields(raw, ',') {
		rep, ok := backRepeatPair(splitTokens(strings.TrimSpace(part)))
		if !ok {
			return nil, false
		}
		out = append(out, rep)
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}

// backRepeatPair folds a keyword list into one X/Y repeat pair.
func backRepeatPair(s []string) (BackRepeat, bool) {
	switch len(s) {
	case 1:
		switch s[0] {
		case "repeat-x":
			return BackRepeat{BackRepeatRepeat, BackRepeatNoRepeat}, true
		case "repeat-y":
			return BackRepeat{BackRepeatNoRepeat, BackRepeatRepeat}, true
		}
		x, ok := backRepeatWord(s[0])
		if !ok {
			return BackRepeat{}, false
		}
		return BackRepeat{x, x}, true
	case 2:
		x, ok := backRepeatWord(s[0])
		if !ok {
			return BackRepeat{}, false
		}
		y, ok := backRepeatWord(s[1])
		if !ok {
			return BackRepeat{}, false
		}
		return BackRepeat{x, y}, true
	}
	return BackRepeat{}, false
}

// isRepeatKeyword reports whether a token belongs to the repeat vocabulary,
// including the repeat-x/repeat-y shorthands.
func isRepeatKeyword(s string) bool {
	if _, ok := backRepeatWord(s); ok {
		return true
	}
	return s == "repeat-x" || s == "repeat-y"
}

// backRepeatWord maps one repeat keyword onto the BackRepeat constants.
func backRepeatWord(s string) (uint8, bool) {
	switch s {
	case "repeat":
		return BackRepeatRepeat, true
	case "no-repeat":
		return BackRepeatNoRepeat, true
	case "space":
		return BackRepeatSpace, true
	case "round":
		return BackRepeatRound, true
	}
	return 0, false
}

// parseBackgroundBox reads a comma-separated background-clip or
// background-origin list of border/padding/content box keywords.
func parseBackgroundBox(raw string) ([]uint8, bool) {
	var out []uint8
	for _, part := range splitFields(raw, ',') {
		switch strings.TrimSpace(part) {
		case "border-box":
			out = append(out, BackBorder)
		case "padding-box":
			out = append(out, BackPadding)
		case "content-box":
			out = append(out, BackContent)
		default:
			return nil, false
		}
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}

// parseBackgroundAttachment reads a comma-separated scroll/fixed/local list.
func parseBackgroundAttachment(raw string) ([]uint8, bool) {
	var out []uint8
	for _, part := range splitFields(raw, ',') {
		switch strings.TrimSpace(part) {
		case "scroll":
			out = append(out, BackAttachScroll)
		case "fixed":
			out = append(out, BackAttachFixed)
		case "local":
			out = append(out, BackAttachLocal)
		default:
			return nil, false
		}
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}

// applyBackground folds a background shorthand into the style: a comma
// separates layers, position and size split at "/", and a colour may only be
// the last layer's. A value that does not fit the grammar drops the whole
// declaration, the way an unrecognised value does.
func applyBackground(st *Style, raw string, ctx Units) bool {
	layers := splitFields(raw, ',')
	if len(layers) == 0 || layers[len(layers)-1] == "" {
		return false
	}
	if len(layers) == 1 && strings.TrimSpace(layers[0]) == "none" {
		st.BackgroundImages = nil
		return true
	}
	var (
		images   []BackImage
		poss     []BackPos
		sizes    []BackSize
		repeats  []BackRepeat
		clips    []uint8
		origins  []uint8
		attaches []uint8
		color    canvas.Color
		colorSet bool
	)
	for idx, layerStr := range layers {
		part := strings.TrimSpace(layerStr)
		head, sizeStr, hasSlash := splitSlash(part)
		var size BackSize
		if hasSlash {
			var ok bool
			if size, ok = parseBackSizeOne(strings.TrimSpace(sizeStr), ctx); !ok {
				return false
			}
		}
		var (
			img     BackImage
			pos     BackPos
			posSet  bool
			repToks []string
			att     uint8
			attSet  bool
			boxes   []uint8
			posToks []string
		)
		for _, tk := range splitTokens(head) {
			switch tk {
			case "border-box":
				boxes = append(boxes, BackBorder)
			case "padding-box":
				boxes = append(boxes, BackPadding)
			case "content-box":
				boxes = append(boxes, BackContent)
			case "scroll":
				att, attSet = BackAttachScroll, true
			case "fixed":
				att, attSet = BackAttachFixed, true
			case "local":
				att, attSet = BackAttachLocal, true
			default:
				if bi, ok := parseBackImage(tk, ctx); ok {
					img = bi
					continue
				}
				if idx == len(layers)-1 {
					if c, err := ParseColor(tk); err == nil {
						color, colorSet = c, true
						continue
					}
				}
				if isRepeatKeyword(tk) {
					repToks = append(repToks, tk)
					continue
				}
				posToks = append(posToks, tk)
			}
		}
		if len(repToks) > 0 {
			rep, ok := backRepeatPair(repToks)
			if !ok {
				return false
			}
			repeats = append(repeats, rep)
		}
		if len(posToks) > 0 {
			p, ok := parseBackPosTokens(posToks, ctx)
			if !ok {
				return false
			}
			pos, posSet = p, true
		}
		sizes = append(sizes, size)
		if posSet {
			poss = append(poss, pos)
		}
		if attSet {
			attaches = append(attaches, att)
		}
		if len(boxes) > 0 {
			origins = append(origins, boxes[0])
			if len(boxes) > 1 {
				clips = append(clips, boxes[1])
			}
		}
		if img.URL != "" || img.Grad != nil {
			images = append(images, img)
		}
	}
	if colorSet {
		st.Background = color
	}
	st.BackgroundImages = images
	st.BackgroundPos = poss
	st.BackgroundSize = sizes
	st.BackgroundRepeat = repeats
	st.BackgroundClip = clips
	st.BackgroundOrigin = origins
	st.BackgroundAttach = attaches
	return true
}
