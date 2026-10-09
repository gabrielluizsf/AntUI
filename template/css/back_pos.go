package css

// BackPos is a background-position: the picture's horizontal and vertical
// anchors. The zero value — both axes anchored to their start edges with no
// offset — is the CSS initial of 0% 0%.
type BackPos [2]BackPosition
