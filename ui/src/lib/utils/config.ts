// src/lib/utils/config.ts

export const baseHttpUrl =
    `http://${import.meta.env.VITE_WS_URL || 'localhost:8080'}`;

export const baseWsUrl =
    `ws://${import.meta.env.VITE_WS_URL || 'localhost:8080'}`;

