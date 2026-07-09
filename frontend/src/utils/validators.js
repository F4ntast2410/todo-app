export function validateEmail(email) {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(String(email).toLowerCase());
}

export function validatePassword(password) {
    // Пароль должен быть минимум 6 символов
    return password && password.length >= 6;
}