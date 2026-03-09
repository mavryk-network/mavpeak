import { formatDistanceStrict, formatDistance } from "date-fns"

export function formatBlockHash(hash: string) {
  if (!hash) {
    return "N/A"
  }
  return `${hash.slice(0, 8)}...${hash.slice(-8)}`
}

export function formatAddress(address: string) {
  if (!address) {
    return "N/A"
  }
  return `${address.slice(0, 8)}...${address.slice(-8)}`
}

function timestampToJSTimestamp(timestamp: number) {
  if (timestamp < 10000000000) {
    return timestamp * 1000
  }
  return timestamp
}

export function formatTimestamp(timestamp: string | number) {
  if (!timestamp) {
    return "N/A"
  }

  if (typeof timestamp === "number") timestamp = timestampToJSTimestamp(timestamp)

  return new Date(timestamp).toLocaleString()
}

export function formatTimestampAgo(timestamp: string | number) {
  if (!timestamp) {
    return "N/A"
  }

  if (typeof timestamp === "number") timestamp = timestampToJSTimestamp(timestamp)

  return formatDistance(new Date(timestamp), new Date())
}

export function formatTimestampAgoStrict(timestamp: string | number) {
  if (!timestamp) {
    return "N/A"
  }

  if (typeof timestamp === "number") timestamp = timestampToJSTimestamp(timestamp)

  return formatDistanceStrict(new Date(timestamp), new Date(), { addSuffix: true })
}

export function formatBalance(mumav: string | bigint | number) {
  if (typeof mumav === "bigint" || typeof mumav === "number") {
    mumav = mumav.toString()
  }
  if (!mumav || typeof mumav !== "string") {
    return "N/A"
  }
  if (mumav === "0") {
    return "0 ꜩ"
  }
  mumav = mumav.padStart(7, "0")
  return `${mumav.substring(0, mumav.length - 6)}.${mumav.substring(mumav.length - 6)} ꜩ`
}

export function formatPercentage(percentage: number | string) {
  return `${Number(percentage).toFixed(2)}%`
}