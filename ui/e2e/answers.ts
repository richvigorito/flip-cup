/**
 * Known question/answer pairs from _.default.yaml
 * Used by game tests to supply correct answers.
 */
export const DEFAULT_QUIZ_ANSWERS: Record<string, string> = {
  'What is the capital of France?': 'Paris',
  'What is 2 + 2?': '4',
  "Who wrote 'Hamlet'?": 'Shakespeare',
  'What is the largest planet in our solar system?': 'Jupiter',
  'In which year did the Titanic sink?': '1912',
  'What is the smallest country in the world?': 'Vatican City',
  'What is the tallest mountain in the world?': 'Mount Everest',
  'Who painted the Mona Lisa?': 'Leonardo da Vinci',
  "Which element has the chemical symbol 'O'?": 'Oxygen',
  'What is the longest river in the world?': 'The Amazon River',
};

/** Fallback: try every known answer until one works. */
export const ALL_ANSWERS = Object.values(DEFAULT_QUIZ_ANSWERS);
