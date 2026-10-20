package template

// cyberStyle paints the Hydra look: near-black glass panels with thin accent
// frames. Hover answers with a faint accent tint and an acc-dim edge, never
// a halo; focus and selection light the frame and bloom a small glow, the
// way the site's HUD does. There are no gradients on the canvas, so the
// glow is faked the honest way — concentric shells of translucent colour.
// Like every style, everything it draws is multiplied by [Scale], so the
// design holds at any resolution.
type cyberStyle struct{}
