// Package css is a small CSS engine for AntUI templates: it reads a stylesheet
// styled after the MDN CSS reference and turns it into the drawing decisions a
// template makes — boxes, borders, radii, colours, typography, shadows and
// layout. It is not a web renderer: it is the part of a CSS that makes sense
// for a canvas of pixels, implemented for templates.
//
// # Selectors
//
// A rule matches by element, class and state. The element name is the widget's
// tag — label, button, checkbox, radio, slider, input (new widgets get tags of
// their own) and body for the window itself. Classes come from the template's
// CSSClasses table; a widget may carry several, and a selector may require any
// of them together, the way .a.b needs both. The universal selector * matches
// every element. Pseudo-classes tie a rule to the interaction state — :hover,
// :focus, :active, :pressed and :checked — so one declaration paints the
// resting state and another paints the same widget under the pointer.
//
//	button                  { background-color: #3E63DD; }
//	button:hover            { background-color: #5472E4; }
//	button.active:hover     { background-color: #2E4BD8; }
//	*                       { box-sizing: border-box; }
//
// # Media queries
//
// An @media rule with min-width, max-width, min-height or max-height (or any
// combination of them) applies only while the window the frame draws at sits
// in that range, which is how a template stays responsive as the window
// resizes: the window is read at every lookup rather than remembered when a
// style is first asked for, so a resize lands on the whole next frame. Three
// more features are read the same way — orientation off the window's two
// edges, resolution off the display's scale in dpi, dpcm or dppx, and
// prefers-color-scheme off the color scheme the system paints in. A condition
// the canvas has no sensor for — a media type, a feature this engine does not
// implement — is reported and treated as satisfied, and so is one whose
// answer the system withheld: a rule is never dropped for want of a
// measurement. Everything else in the rule follows the normal cascade.
//
//	@media (min-width: 720px) { button { font-size: 20px; } }
//	@media (max-height: 480px) { body { font-size: 14px; } }
//	@media (orientation: landscape) { nav { flex-direction: row; } }
//	@media (prefers-color-scheme: dark) { body { background-color: #101014; } }
//
// # Feature tests
//
// An @supports rule gates its body on what the engine itself can do. The
// condition is read against the same cascade the stylesheet runs on, so a test
// passes only for a property the engine has holding a value it reads, or a
// selector it can evaluate — and selector() asks the very reader that runs
// every rule in the sheet. and, or, not and the groups around them read as
// they read in CSS. A test that fails leaves its block out of the sheet, which
// is what the test is for; a condition nothing can make sense of is reported
// before the block goes, so a half-written condition is heard rather than
// silently obeyed. The answer rests on the engine alone — it is read while the
// sheet parses, not measured against the window — so the same stylesheet takes
// the same branches in every frame. font-format() answers yes for TrueType,
// the outlines an @font-face reads, and no for the woff containers.
// font-tech() answers yes for color-CBDT, the colour bitmaps an @font-face
// reads and every frame draws, and no for the rest.
//
//	@supports (display: grid) { .card { display: grid; } }
//	@supports not (display: subgrid) { .card { display: grid; } }
//	@supports selector(button:hover) and (color: #fff) { button { color: #fff; } }
//
// # Fonts
//
// An @font-face block reads a TrueType file for a family while the sheet
// parses, at the size the program draws its own text at, and the family,
// font-weight and font-style it declares are how a font-family asking for
// that name finds it: the leaning face first when the style leans and the
// straight one first when it does not, and inside that the weight asked for
// and then the nearest ones the spec turns up. A source this engine does not
// read — local(), naming a font this machine already has, or a woff container
// — is reported and passed over, and so is a file that will not open: a
// @font-face leaves its family with no face rather than failing the sheet that
// held it.
//
//	font-family: "Roboto", sans-serif;
//	@font-face {
//	    font-family: "Roboto";
//	    src: url("Roboto-Regular.ttf") format("truetype");
//	    font-weight: 400;
//	}
//
// A font-family list is not only a chooser of one file: the names after the
// first are the faces a rune the first one lacks falls through to, in order,
// ending at the face the program draws with by default — so a character one
// font holds and another does not still comes out in the right shape. The
// chain is drawn at the size the rest of the program's text is, so naming a
// family changes which letters draw and not how big they come out; the
// stylesheet's own font-size is what sizes them. A face that already leans or
// already carries the weight asked for has no synthetic italic or bold drawn
// over it, and a line is as tall as the face in it, so line-height multiplies
// that height rather than the built-in font's. A list naming no family the
// sheet holds a file for answers with no face, and the canvas keeps drawing in
// the face it already had.
//
// A file whose glyphs are pictures rather than outlines — the CBDT colour
// bitmaps an emoji font is made of — is read beside the rest: the strike the
// file cut nearest the size being drawn supplies the picture, and its
// bearings and advance are what put it on the line and move the pen on by.
// The chain above decides whose rune draws as it always did — a rune the
// first face holds a glyph for comes from that file, whether it draws an
// outline or a picture of it, and only the runes it lacks reach the faces
// after it — and a picture leans and doubles under font-style and font-weight
// the way an outline does.
//
// # Values
//
// Lengths are px (scaled with the window, like everything else a template
// draws) or percentages of the window width. Colors may be named, hex, or
// rgb()/rgba(). Border styles cover none, solid, dashed, dotted and double.
// The box model supports width, height, min/max constraints, margin, padding
// and box-sizing; the border model supports border, border-width/style/color
// and border-radius. background-color and background paint the surface, and
// background-image layers url() images registered with the template or the
// linear/radial/conic gradient functions, with background, position, size,
// repeat, clip, origin and attachment controlling where each layer sits and
// how it clips. box-shadow and text-shadow paint shadows (a list, with inset,
// blur and spread, the colour defaulting to currentColor), outline and
// outline-offset ring the box outside it, and filter/backdrop-filter apply
// grayscale, sepia, invert, brightness, contrast, hue-rotate, blur and
// drop-shadow to the box or what is behind it. transform turns a widget
// through translate/scale/rotate/skew/matrix lists around its
// transform-origin, a pivot with keyword, length or percentage values whose
// initial value sits at the box's centre. opacity fades, and font-size,
// font-family, line-height, text-align and color drive the text.
// display:none leaves the widget
// out of the frame entirely, and position:absolute with top/left lifts a
// widget out of the flow so it can be parked against the window edges.
//
// Transitions and animations move those numbers between frames. A transition
// eases a property from its resting value to its hovered (or focused, or
// active, or checked) one over a duration after an optional delay:
//
//	button { transition: background-color 0.2s ease-out; }
//	button:hover { background-color: #5472E4; }
//
// The shorthand (transition) and its longhands (transition-property,
// transition-duration, transition-timing-function, transition-delay) both
// work, and their comma-separated lists pair up the way CSS pairs them. An
// @keyframes block names a series of frames — from, to, or percentages —
// and the animation shorthand runs them, with a duration, an easing, a
// delay, an iteration count, a direction (reverse/alternate/alternate-reverse)
// and a fill-mode (backwards/forwards/both) pinning the frame before the
// animation starts or keeps it after it ends:
//
//	@keyframes pulse {
//		from { opacity: 0.4; }
//		to   { opacity: 1; }
//	}
//	button { animation: pulse 1s ease-in-out infinite alternate; }
//
// The animatable properties are the colours, opacity, the length properties
// (widths, heights, margins, paddings, insets, font-size, letter/word-spacing,
// line-height), transform, transform-origin, the shadows and the filters.
// Everything else snaps at the swap; border widths do not interpolate. The
// engine draws the motion itself, carried by the template's frame clock, so
// a hover polish needs no timers from the developer.
//
// Typography: font-weight (fake bold above 600, bolder and lighter
// approximate), font-style italic/oblique (skewed strokes), line-height,
// letter-spacing, word-spacing, text-transform (uppercase/lowercase/capitalize),
// text-decoration (underline/overline/line-through), white-space
// (normal/nowrap/pre/pre-wrap/pre-line), overflow-wrap (break-word/anywhere),
// text-overflow (ellipsis), text-align (left/center/right/justify) and
// vertical-align (sub/super/middle/top/bottom/text-top/text-bottom and
// lengths). They inherit where CSS inherits them and draw through the canvas
// text style.
//
// The full MDN CSS reference (https://developer.mozilla.org/en-US/docs/Web/CSS/Reference)
// names far more than any canvas can paint; the engine implements the surface
// that maps directly onto it, and every property it does not know is ignored
// with a warning rather than an error.
//
// # Layout containers
//
// Two properties hand a block over to a layout engine instead of the vertical
// flow: display:flex makes it a flex container, display:grid a grid container.
// Both take the whole box — margin, border, padding, width, height and
// min/max — and place the widgets drawn inside them, which the template calls
// CSS.Flex and CSS.Grid. The gap, alignment, order and sizing properties those
// two engines read are:
//
//	flex-direction, flex-wrap, flex-grow, flex-shrink, flex-basis, order,
//	gap, row-gap, column-gap, justify-content, justify-items, justify-self,
//	align-content, align-items, align-self
//
// A flex item is sized along the main axis by its basis, grown and shrunk
// against what is left, and placed on the cross axis by align-items. The main
// axis runs with flex-direction, wraps with flex-wrap, and justify-content
// packs the items along it.
//
// A grid container is split into tracks by grid-template-columns and
// grid-template-rows. A track is a length, auto, a fraction (fr), minmax(),
// the content keywords min-content, max-content and fit-content(), or
// repeat() — with a count, or with auto-fill and auto-fit to work out the
// count from the room the container has. gap separates the tracks,
// grid-template-areas names the rectangles of cells, and an item places itself
// with grid-column, grid-row and grid-area, by line number, by span, or by the
// names given to the lines in a template, counting from the end with a leading
// minus. grid-auto-flow chooses whether items fill the axis in turn (row) or
// the other one (column), sparse or dense, and order moves an item along
// without moving it in the source.
//
//	grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(120px, 1fr)); gap: 8px; }
//	grid { grid-template-columns: [full-start] 1fr [main-start] 2fr [main-end] 1fr [full-end]; }
//	item { grid-column: main-start / main-end; }
//
// subgrid is not supported: asking a nested container to adopt the lines of
// the one around it is an error, not a silent approximation. A grid inside
// another grid, or inside a flex, lays its own tracks out from scratch inside
// the box it was given.
//
// # Columns and tables
//
// A block the stylesheet marks with the multicol class becomes a multi-column
// container, which the template calls CSS.MultiCol, and the content inside it is
// split into equal columns with a gutter between them. The properties that
// decide the split are:
//
//	column-count, column-width, column-gap, column-rule, column-rule-width,
//	column-rule-style, column-rule-color, column-fill, columns,
//	break-inside, page-break-inside
//
// column-count says how many columns there are, column-width how wide one of
// them is (which is also a way of asking for as many as fit), and the smaller of
// the two wins. column-gap is the same gap flex and grid read, and column-rule —
// a border's three declarations, which take a colour like currentColor — is
// painted down the middle of it. column-fill: auto spends the whole height on one
// column before starting the next; the default balance spreads the height over
// them instead, and a container with no height of its own is as tall as the
// tallest of its columns. A container with a height has one to cut a box
// against: a box that runs out of room is cut and continued in the next column,
// and a box with break-inside: avoid is moved whole to the next one instead. What
// does not fit the columns the container asked for hangs below the last one.
//
//	multicolumn { column-count: 3; column-gap: 16px; column-rule: 1px solid #ddd; }
//	quote     { break-inside: avoid; }
//
// A box is cut by giving its fragments a draw pass each, and the callback is run
// once per pass, so it must draw the same widgets in the same order every run. A
// box that would need more than 64 fragments is drawn whole in the piece that
// overflows rather than costing the frame more passes.
//
//	display:table makes a block a table container, which the template calls
//	CSS.Table, and the widgets inside it are the cells of the rows the callback
//	draws with CSS.TableRow. A column is as wide as the widest cell in it, a row
//	as tall as the tallest cell in it, and a table with a width of its own shares
//	that width out over the columns in proportion to what they asked for — a wide
//	column keeps more of the extra room than a narrow one. A cell's content is
//	placed in its row by the vertical-align of the cell's own style. The interior
//	display values are read, not folded into display:table the way a browser folds
//	them, so a stylesheet can say what a widget is:
//
//	table  { display: table; }
//	cell   { display: table-cell; vertical-align: middle; }
//	caption{ display: table-caption; }
//
// A widget drawn straight in a table, without a row around it, becomes a row of
// one cell, which is what a browser's anonymous boxes do with a child that did
// not ask to be a row. A widget marked display:table-caption is not a cell: it
// is painted above the table, across its whole width.
//
// Tables are laid out with separated borders, the one case of the two that
// collapses nothing: border-collapse and border-spacing are read as ordinary
// properties and change nothing about where a border is painted. A table in a
// cell, a column, a flex item or a grid item is an item of it — as wide as its
// own columns ask for, and centred in what it was given — and a column block in
// any of them keeps its own columns, re-solved at the box it really ended up in.
//
// float and clear are read and warned about but change nothing: this flow has no
// line box to float a box inside, and nothing here floats for clear to push
// past. break-before and break-after are not supported; only break-inside, which
// is what a multi-column cut needs.
package css
