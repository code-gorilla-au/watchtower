import { formatDistanceToNow } from "date-fns";

/**
 * Formats the given date to display the relative time (distance to now).
 */
export function formatDate(date: Date) {
	return formatDistanceToNow(date);
}

/**
 * Truncates the provided string to a maximum of 100 characters and appends an ellipsis ("...") at the end.
 */
export function truncate(str: string) {
	const maxLength = 100;

	if (str.length <= maxLength) {
		return str;
	}

	const sub = str.substring(0, maxLength);
	return `${sub}...`;
}

/**
 * Converts a string from camelCase, snake_case, or PascalCase to `Sentence case`.
 */
export function toSentenceCase(str: string): string {
	if (!str) {
		return str;
	}

	const result = str
		.replace(/[_-]+/g, " ")
		.replace(/([a-z])([A-Z])/g, "$1 $2")
		.replace(/([A-Z])([A-Z][a-z])/g, "$1 $2")
		.trim()
		.toLowerCase();

	return result.charAt(0).toUpperCase() + result.slice(1);
}
