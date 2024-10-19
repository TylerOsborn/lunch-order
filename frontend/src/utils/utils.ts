function zeroPad(num: number, places: number) {
  return String(num).padStart(places, '0');
}

export function mondayDate() {
  const today = new Date();
  const day = today.getDay();
  const diff = today.getDate() - day + (day == 0 ? -6 : 1);
  const monday = new Date(today.setDate(diff));
  const year = zeroPad(monday.getFullYear(), 4);
  const month = zeroPad(monday.getMonth() + 1, 2);
  const date = zeroPad(monday.getDate(), 2);
  return `${year}-${month}-${date}`;
}

export function thursdayDate() {
  const today = new Date();
  const day = today.getDay();
  const diff = today.getDate() - day + (day == 0 ? -6 : 4);
  const monday = new Date(today.setDate(diff));
  const year = zeroPad(monday.getFullYear(), 4);
  const month = zeroPad(monday.getMonth() + 1, 2);
  const date = zeroPad(monday.getDate(), 2);
  return `${year}-${month}-${date}`;
}

export function getTodayDate() {
  const today = new Date();
  const year = zeroPad(today.getFullYear(), 4);
  const month = zeroPad(today.getMonth() + 1, 2);
  const date = zeroPad(today.getDate(), 2);
  return `${year}-${month}-${date}`;
}

export function getNameFromLocalStorage(): string {
  return localStorage.getItem('name') || '';
}

export function setNameToLocalStorage(name: string) {
  localStorage.setItem('name', name);
}

export function getUUIDFromLocalStorage(): string {
  return localStorage.getItem('uuid') || '';
}

export function setUUIDToLocalStorage(uuid: string) {
  localStorage.setItem('uuid', uuid);
}
