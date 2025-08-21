import { formatDistanceToNow } from "date-fns";

export function formatDate(date: Date) {
	return formatDistanceToNow(date);
}

/**
 *
 * @param str
 */
export function truncate(str: string) {
	const sub = str.substring(0, 100);
	return `${sub}...`;
}
