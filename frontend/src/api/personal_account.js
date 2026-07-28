import { apiRequest } from './client.js';

export async function changeInfo(payload) {
    return apiRequest('/auth/me', {
        method: 'PUT',
        body: JSON.stringify(payload)
    });
}
export async function changeEmailVerify(payload) {
    return apiRequest('/auth/me/email/verify', {
        method: 'PUT',
        body: JSON.stringify(payload)
    });
}

export async function changePassword(payload) {
    return apiRequest('/auth/me/password', {
        method: 'PUT',
        body: JSON.stringify(payload)
    });
}

export async function checkLinkTg() {
    return apiRequest('/auth/telegram/link', {
        method: 'GET',
    });
}

export async function linkTg(user) {
    return apiRequest('/auth/telegram/link', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(user)
    });
}
export async function setEmail(email) {
    return apiRequest('/auth/me/link/email', {
        method: 'POST',
        body: JSON.stringify({ email: email })
    });
}
export async function setEmailVerify2FA(email, inputCode) {
    return apiRequest('/auth/me/link/verify2fa', {
        method: 'POST',
        body: JSON.stringify({ email: email, input_code: inputCode })
    });
}
export async function setEmailPassword(email, password) {
    return apiRequest('/auth/me/link/password', {
        method: 'POST',
        body: JSON.stringify({ email, password })
    });
}