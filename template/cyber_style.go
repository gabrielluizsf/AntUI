package template

// cyberStyle paints the futuristic look: dark glass panels, thin neon edges,
// and a soft glow around whatever the pointer is over. There are no gradients
// on the canvas, so the glow is faked the honest way — concentric shells of
// translucent colour, the same way the glow sits on real neon. Like every
// style, everything it draws is multiplied by [Scale], so the design holds at
// any resolution.
type cyberStyle struct{}
