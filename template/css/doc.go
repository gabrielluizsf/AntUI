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
// An @media rule with min-width or max-width (or both) applies only while the
// window's width is in that range, which is how a template stays responsive as
// the window resizes. Everything else in the rule follows the normal cascade.
//
//	@media (min-width: 720px) { button { font-size: 20px; } }
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
// initial value sits at the box's centre. opacity fades, and
// font-size/text-align/color drive the text. display:none leaves the widget
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
// float and clear are read and warned about but change nothing: this flow has no
// line box to float a box inside, and nothing here floats for clear to push
// past. break-before and break-after are not supported; only break-inside, which
// is what a multi-column cut needs.
package css
