package svg

import "github.com/gabrielluizsf/antui/canvas"

// Style is the painting a node was given: how thick its stroke is, how its ends
// and corners are treated, how transparent it is, and whether it is drawn at
// all. It is kept whole on the node so a caller that wants to restyle a drawing
// can see what the file asked for before deciding what to change.
type Style struct {
	Fill      canvas.Color
	HasFill   bool
	FillRule  canvas.FillRule
	Stroke    canvas.Color
	HasStroke bool
	// ClipRule is the rule the shapes of a clipPath are read with when this
	// element is the one they cut — the clip-rule of the clipPath, inherited
	// down to the shapes inside it the way fill-rule is — and it says which of
	// the insides of a shape that crosses itself counts.
	ClipRule canvas.FillRule
	// FillCurrent and StrokeCurrent say that the paint was the keyword
	// `currentColor`, so the colour the node was given is not the one it is
	// drawn in: it takes the colour the painting was asked for. One drawing
	// written once therefore paints in any colour, which is how a single icon
	// serves a button in every state without being written again.
	FillCurrent   bool
	StrokeCurrent bool
	// FillContext and StrokeContext say that the paint was the keyword
	// `context-fill` or `context-stroke`: the colour of the content of a marker
	// comes from the fill or the stroke of the element that named the marker.
	// The colour is only known where the marker is drawn, and one marker's
	// picture is shared by every shape that names it, so the keyword is kept as
	// a flag and resolved in [palette] at paint time rather than here.
	FillContext   bool
	StrokeContext bool
	Width         float64
	Cap           canvas.LineCap
	Join          canvas.LineJoin
	MiterLimit    float64
	// VectorEffect is `vector-effect`: the ways the element's own user
	// coordinate system is held still against the host's transformations,
	// including the stroke that does not scale — a hairline that stays a
	// hairline however much the drawing is zoomed. It is not inherited, as the
	// spec has it, so [Style.with] takes the whole of it off again before
	// reading what the element itself said, the same as it does for pathLength,
	// while the coordinate system the effect makes of the element is what the
	// shapes inside it are drawn in. See [vectorEffectMatrix].
	VectorEffect  VectorEffect
	Opacity       float64
	FillOpacity   float64
	StrokeOpacity float64
	// Dash is the stroke broken into pieces along the path, with every length
	// of it — the dashes, the gaps and the offset into them — written in the
	// drawing's own units. On a shape whose element wrote a pathLength it
	// holds those pieces stretched to the length that was declared, which is
	// what the shape is cut with; what the shapes around it inherit is the
	// pattern as it was written, since the length it was stretched against
	// belongs to this shape alone. The cut itself happens where the shape is
	// built and not where it is painted — see [Image.build] — so the pattern
	// reaches the pixels the way everything else the drawing wrote does.
	Dash         *canvas.Dash
	Hidden       bool
	Transform    canvas.Matrix
	HasTransform bool
	// FontSize is how tall the writing is, in the drawing's own units, and
	// Anchor is where along the pen a piece of it hangs. Both only matter to a
	// `<text>`, but both are inherited like the rest, so a size written on a
	// `<g>` is the size everything inside it is written at.
	//
	// A FontSize of zero is no size written at all rather than no writing: the
	// reader of the whole drawing decides what that means, which is what lets
	// a size written once on the root govern every text under it. See
	// [readFontSize].
	FontSize float64
	Anchor   TextAnchor
	// PreserveSpace is `xml:space="preserve"`: the writing is drawn exactly as
	// it was written in the file, whitespace runs and line breaks and all,
	// instead of being squeezed onto one line the way a file wrapped across
	// lines is normally read. Like the size and the anchor above it, it is
	// inherited, and `xml:space="default"` on a nearer element takes it off
	// again.
	PreserveSpace bool
	// base is where the drawing sits, for an address inside it to resolve
	// against: the file the drawing was read from, or the `xml:base` an element
	// wrote, which is inherited like `xml:space` above. It is empty when the
	// drawing came from no file and said nothing of its own, and then a
	// relative address is left exactly as it was written. See [Image.fillPicture].
	base string
	// FontFamily, FontWeight and LetterSpacing are the rest of what says how
	// the writing is set: the faces it asks for, how heavy it is, and how much
	// room is left after every letter. All three are inherited like the size
	// and the anchor above, and all three only mean anything to writing.
	//
	// The family is kept the way the CSS side keeps one — the list as it was
	// written, so its names can be asked for one at a time — and the answer
	// comes from the canvas chain: a drawing read here is a single file with no
	// `@font-face` behind it, so no name in any list is one that resolves to a
	// file of its own, and the writing comes out in the face the program draws
	// with by default, exactly as CSS draws text whose family the sheet holds
	// no face for. See [Style.with].
	FontFamily string
	// FontWeight is one of the hundreds CSS matches weights on, or zero when
	// nothing said. `bold` and a number of the scale are what they say;
	// `bolder` and `lighter` are the step CSS takes from the weight in hand.
	// Only the heavy end of the scale draws differently — the canvas fakes a
	// weight a face does not carry by drawing it a second time, and there is
	// no second pass that takes ink away.
	FontWeight uint16
	// LetterSpacing is extra room after every letter, in the drawing's own
	// units. Zero is `letter-spacing="normal"` and also no letter-spacing at
	// all, which are the same thing.
	LetterSpacing float64
	// A fill or a stroke may be a gradient or a pattern rather than a colour,
	// and the three never both hold: `fillGradient` is set only when the paint
	// named a gradient that was in the drawing, `fillPattern` the same for a
	// pattern, and `fillRef` is the id it named until the drawing has been read
	// far enough to say what that is. The fallback colour is the one written
	// after the `url(...)`, which is what the shape is painted with when the
	// reference is not there at all.
	fillGradient          *gradient
	fillPattern           *pattern
	fillRef               string
	fillFallback          canvas.Color
	fillFallbackCurrent   bool
	hasFillFallback       bool
	strokeGradient        *gradient
	strokePattern         *pattern
	strokeRef             string
	strokeFallback        canvas.Color
	strokeFallbackCurrent bool
	hasStrokeFallback     bool
	// clipRef is the id a `clip-path="url(#id)"` named, kept the same way a
	// fill reference is: only the whole drawing can say what it is, and the
	// clip it turns into is not inherited by the children — they are drawn
	// into this element's picture first and the picture is cut once. It is
	// cleared as it is spent, see [Image.resolveClip].
	clipRef string
	// maskRef is the id a `mask="url(#id)"` named, kept the same way a clip
	// reference is, and cleared as early: what it turns into is not inherited
	// either, so it is taken off the style before the children are built from
	// it and followed once the whole of this element is there — the region a
	// mask reaches is measured against everything it paints. See
	// [Image.maskNamed] and [maskUnder].
	maskRef string
	// markerStart, markerMid and markerEnd are the `<marker>`s this element
	// draws at the first, the in-between and the last vertex of its shape, and
	// they are inherited, which is what sets them apart from the clip, the
	// mask and the filter above: a marker is not a cut or a picture laid over
	// one shape, it is a mark on a vertex, and a `<g marker-end>` is what every
	// shape under it puts on its own last vertex. So the references turn into
	// the markers themselves rather than into places on the node, and the
	// children take them along with the fill. See [Image.resolveMarkers].
	markerStart *markerDef
	markerMid   *markerDef
	markerEnd   *markerDef
	// markerStartRef, markerMidRef and markerEndRef are those references — the
	// ids the three `url(#id)`s named — until the whole drawing has been read
	// far enough to say what they are, and each is cleared as it is spent so
	// that a group says what it has to say once and not once for every shape
	// inside it.
	markerStartRef string
	markerMidRef   string
	markerEndRef   string
	// filters is the list a `filter="blur(1) grayscale(1)"` gave the element,
	// read where it is written rather than once the whole drawing is in hand:
	// none of the functions needs to know what any id in the file is, and what
	// is not read is said at once. It is not inherited either — a `<g filter>`
	// turns the one picture it makes of its children once — so it comes off
	// the style before the children are built from it and lives on the node;
	// see [readFilters] and [applyFilters].
	filters []filterOp
	// pathLength is the total length the drawing says the path on this element
	// has, in the drawing's own units, which is what the stroke of this shape
	// is measured out against instead of the geometry — see [applyPathLength].
	// It is not inherited, for the same reason the clip, the mask and the
	// filter are not: a length declared by a group is a length of the group's
	// own shape, and a group has none. [Style.with] takes it off again for
	// every element that does not write one of its own.
	pathLength float64
}
