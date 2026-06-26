// TypeScript port of shared/timecalc/calculate.go for live preview in the UI.
// The authoritative calculation still happens server-side on submit.

function round(value: number, places: number): number {
  const factor = Math.pow(10, places);
  return Math.round(value * factor) / factor;
}

// Returns minutes since midnight, or null if the time string is invalid.
function parseTime(raw: string): number | null {
  const timeStr = raw.toLowerCase().trim();
  if (!timeStr) return null;

  let hours: number;
  let minutes: number;

  if (timeStr.includes(":")) {
    if (timeStr.length < 6 || timeStr.length > 7) return null;
    const parts = timeStr.split(":");
    hours = parseInt(parts[0], 10);
    if (isNaN(hours) || hours < 1 || hours > 12) return null;

    let minutesPart = parts[1];
    if (minutesPart.endsWith("am")) {
      if (hours === 12) hours = 0;
      minutesPart = minutesPart.slice(0, -2);
    } else if (minutesPart.endsWith("pm")) {
      if (hours !== 12) hours += 12;
      minutesPart = minutesPart.slice(0, -2);
    } else {
      return null;
    }

    minutes = parseInt(minutesPart, 10);
    if (isNaN(minutes) || minutes < 0 || minutes > 59) return null;
  } else {
    if (timeStr.length < 3 || timeStr.length > 4) return null;
    minutes = 0;
    let numStr: string;
    if (timeStr.endsWith("am")) {
      numStr = timeStr.slice(0, -2);
      hours = parseInt(numStr, 10);
      if (hours === 12) hours = 0;
    } else if (timeStr.endsWith("pm")) {
      numStr = timeStr.slice(0, -2);
      hours = parseInt(numStr, 10);
      if (hours !== 12) hours += 12;
    } else {
      return null;
    }
    if (isNaN(hours) || hours < 1 || hours > 24) return null;
  }

  return hours * 60 + minutes;
}

export interface CalcResult {
  hours: number | null;
  error: string | null;
}

export function calculateTime(
  start: string,
  end: string,
  lunch: string,
  extra: string
): CalcResult {
  const startMin = parseTime(start);
  if (startMin === null) {
    return { hours: null, error: "Enter a valid start time (e.g. 9:00am)" };
  }
  const endMin = parseTime(end);
  if (endMin === null) {
    return { hours: null, error: "Enter a valid end time (e.g. 5:00pm)" };
  }
  if (endMin < startMin) {
    return { hours: null, error: "End time is before start time" };
  }

  const subtract = lunch.trim() === "" ? 0 : parseFloat(lunch);
  if (isNaN(subtract)) {
    return { hours: null, error: "Lunch must be a number" };
  }
  const add = extra.trim() === "" ? 0 : parseFloat(extra);
  if (isNaN(add)) {
    return { hours: null, error: "Extra must be a number" };
  }

  const totalHours = (endMin - startMin) / 60 - subtract + add;
  return { hours: round(totalHours, 2), error: null };
}
