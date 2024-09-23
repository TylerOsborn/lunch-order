function zeroPad(num: number, places: number) {
    return String(num).padStart(places, '0');
}

export function mondayDate() {
    const today = new Date();
    const day = today.getDay();
    const diff = today.getDate() - day + (day == 0 ? -6 : 1);
    const monday = new Date(today.setDate(diff))
    const year = zeroPad(monday.getFullYear(), 4);
    const month =zeroPad(monday.getMonth() + 1, 2);
    const date = zeroPad(monday.getDate(), 2);
    return `${year}-${month}-${date}`
}

export function thursdayDate() {
    const today = new Date();
    const day = today.getDay();
    const diff = today.getDate() - day + (day == 0 ? -6 : 4);
    const monday = new Date(today.setDate(diff))
    const year = zeroPad(monday.getFullYear(), 4);
    const month = zeroPad(monday.getMonth() + 1, 2);
    const date = zeroPad(monday.getDate(), 2);
    return `${year}-${month}-${date}`
}

export function setNameCookie(name: string) {
    const d = new Date();
    d.setTime(d.getTime() + (365*24*60*60*1000)); // 1 year expiration
    const expires = "expires=" + d.toUTCString();
    document.cookie = "username=" + name + ";" + expires + ";path=/";
}

export function getNameFromCookie(): string {
    const name = "username=";
    const decodedCookie = decodeURIComponent(document.cookie);
    const ca = decodedCookie.split(';');
    for(let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }

    return '';
}