export function debouncedInput(fn: (input: string) => void, delay: number = 500) {
	let timeout: number;

	return (input: string) => {
		window.clearTimeout(timeout);
		timeout = window.setTimeout(() => {
			fn(input);
		}, delay);
	};
}
