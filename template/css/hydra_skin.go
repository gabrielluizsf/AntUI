package css

// hydraSkin is the Hydra look — the Solo Task site's HUD, expressed in the
// template's own tags. It is applied on top of the program's stylesheet, so
// it repaints the widgets without touching the file's layout.
const hydraSkin = `
body { background-color: #04070F; color: #FFFFFF; }
label { color: #FFFFFF; }

button {
	background-color: #57C1FF;
	color: #04070F;
	border-width: 0;
	clip-path: polygon(0 0, calc(100% - 10px) 0, 100% 10px, 100% 100%, 10px 100%, 0 calc(100% - 10px));
	font-family: "JetBrains Mono";
	font-weight: 700;
	letter-spacing: 2px;
	text-transform: uppercase;
}
button:hover {
	background-color: #FFFFFF;
	color: #04070F;
	box-shadow: 0 0 20px rgba(87, 193, 255, 0.6);
}
button:active {
	background-color: #57C1FF;
	color: #04070F;
	box-shadow: 0 0 30px rgba(87, 193, 255, 0.7);
}

input, select, textarea {
	background-color: #020306;
	border-color: #0E2233;
	border-width: 1px;
	color: #FFFFFF;
	clip-path: polygon(6px 0, calc(100% - 6px) 0, 100% 6px, 100% calc(100% - 6px), calc(100% - 6px) 100%, 6px 100%, 0 calc(100% - 6px), 0 6px);
}
input:hover, select:hover, textarea:hover { border-color: #224B6A; }
input:focus, select:focus, textarea:focus {
	border-color: #57C1FF;
	box-shadow: 0 0 15px rgba(87, 193, 255, 0.25);
}

checkbox, radio, switch { border-color: #224B6A; }
checkbox:checked, radio:checked, switch:checked { border-color: #57C1FF; }
checkbox:focus, radio:focus, switch:focus { box-shadow: 0 0 12px rgba(87, 193, 255, 0.4); }

slider { background-color: #57C1FF; border-color: #224B6A; }
progress { background-color: #57C1FF; border-color: rgba(255, 255, 255, 0.05); }

datepicker {
	background-color: #050B19;
	border-color: rgba(87, 193, 255, 0.35);
	clip-path: polygon(12px 0, calc(100% - 12px) 0, 100% 12px, 100% calc(100% - 12px), calc(100% - 12px) 100%, 12px 100%, 0 calc(100% - 12px), 0 12px);
}
`
