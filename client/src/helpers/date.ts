// Output: "Wed Jun 20, 2023 - 7:38 PM"
export const formatTimestamp = (timestamp: string): string => {
  const options: Intl.DateTimeFormatOptions = {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: 'numeric',
    hour12: true,
  }

  const date = new Date(timestamp)
  const currentYear = new Date().getFullYear()
  const formattedDate = new Intl.DateTimeFormat(undefined, options).format(date)

  if (date.getFullYear() === currentYear) {
    const withoutYear = formattedDate.replace(`, ${currentYear}`, '')
    return withoutYear
  }

  return formattedDate
}
