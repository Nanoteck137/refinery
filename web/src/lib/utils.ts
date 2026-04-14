import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function capitilize(s: string) {
  if (s.length === 0) return "";
  return s[0].toUpperCase() + s.substring(1);
}

export function formatTime(s: number) {
  const min = Math.floor(s / 60);
  const sec = Math.floor(s % 60);

  return `${min}:${sec.toString().padStart(2, "0")}`;
}

export function formatDuration(ms: number): string {
  if (ms < 1000) return `${ms}ms`;

  const s = Math.floor(ms / 1000);
  const subsec = ms % 1000;
  const h = Math.floor(s / 3600);
  const m = Math.floor((s % 3600) / 60);
  const sec = s % 60;

  return [
    h && `${h}h`,
    m && `${m}m`,
    (sec || subsec) && `${sec}.${String(subsec).padStart(3, "0")}s`,
  ]
    .filter(Boolean)
    .join(" ");
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function isRoleAdmin(role: string) {
  switch (role) {
    case "super_user":
    case "admin":
      return true;
    default:
      return false;
  }
}

export function getPagedQueryOptions(searchParams: URLSearchParams) {
  const query: Record<string, string> = {};

  const page = searchParams.get("page");
  if (page) {
    query["page"] = page;
  }

  return query;
}

export function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
