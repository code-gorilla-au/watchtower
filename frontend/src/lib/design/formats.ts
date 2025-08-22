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
