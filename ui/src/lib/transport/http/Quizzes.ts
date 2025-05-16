import type { QuestionSet } from '$lib/types/QuestionSet'; 
import { baseHttpUrl } from '$lib/utils/config';



export async function fetchQuizzes(): Promise<QuestionSet[]> {
    try {
        const httpUrl = `${baseHttpUrl}/quizzes`;

        const res = await fetch(httpUrl);
        const quizzes: QuestionSet[] = await res.json()

        console.log('✅ Fetched quizzes:', quizzes);
        return  quizzes;
    } catch (err) {
        console.error('Failed to fetch question sets:', err);
        return [];
    }
}
