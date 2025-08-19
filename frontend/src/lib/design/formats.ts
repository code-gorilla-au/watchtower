import { formatDistanceToNow } from "date-fns";

export function formatDate(date: Date) {
	return formatDistanceToNow(date);
}
