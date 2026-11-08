/**
 * @param saturation 0 <= saturation <= 1
 * @param value 0 <= value <= 1
 */
export function generateRandomNonYellowHue() {
	// exclude 45° to 75° (yellow), range of 30°
	return (Math.round(Math.random() * 330) + 75) % 360;
}

/** https://www.rapidtables.com/convert/color/hsl-to-rgb.html
	* @param hue: 0 <= hue < 360
	* @param saturation: 0 <= saturation <= 1
	* @param lightness: 0 <= lightness <= 1
	*/
export function hslToRgb(hue: number, saturation: number, lightness: number) {
	const c = (1 - Math.abs(2 * lightness - 1)) * saturation
	const x = c * (1 - Math.abs((hue / 60) % 2 - 1))
	const m = lightness - c / 2

	let r, g, b
	if (hue < 60) {
		r = c
		g = x
		b = 0
	} else if (hue < 120) {
		r = x
		g = c
		b = 0
	} else if (hue < 180) {
		r = 0
		g = c
		b = x
	} else if (hue < 240) {
		r = 0
		g = x
		b = c
	} else if (hue < 300) {
		r = x
		g = 0
		b = c
	} else if (hue < 360) {
		r = c
		g = 0
		b = x
	} else {
		r = -10000;
		g = -10000;
		b = -10000;
	}

	r = Math.round((r + m) * 255);
	g = Math.round((g + m) * 255);
	b = Math.round((b + m) * 255);

	return [r, g, b];
}

/** https://www.rapidtables.com/convert/color/hsv-to-rgb.html
	* @param h: 0 <= h < 360
	* @param s: 0 <= s <= 1
	* @param v: 0 <= v <= 1
	*/
function hsvToRgb(h: number, s: number, v: number) {
	const c = v * s;
	const x = c * (1 - Math.abs((h / 60) % 2 - 1));
	const m = v - c;

	let r, g, b
	if (h < 60) {
		r = c
		g = x
		b = 0
	} else if (h < 120) {
		r = x
		g = c
		b = 0
	} else if (h < 180) {
		r = 0
		g = c
		b = x
	} else if (h < 240) {
		r = 0
		g = x
		b = c
	} else if (h < 300) {
		r = x
		g = 0
		b = c
	} else if (h < 360) {
		r = c
		g = 0
		b = x
	} else {
		r = -10000;
		g = -10000;
		b = -10000;
	}

	r = Math.round((r + m) * 255);
	g = Math.round((g + m) * 255);
	b = Math.round((b + m) * 255);

	return [r, g, b];
}

/**
	* @param r 0<=r<=255
	* @param g 0<=r<=255
	* @param b 0<=r<=255
	*/
export function rgbToHex(r: number, g: number, b: number): string {
	return `#${r.toString(16)}${g.toString(16)}${b.toString(16)}`
}
