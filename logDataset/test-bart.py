from transformers import pipeline
from transformers import AutoTokenizer, TFAutoModelForSeq2SeqLM

model = TFAutoModelForSeq2SeqLM.from_pretrained('VidhuMathur/bart-log-summarization')

tokenizer = AutoTokenizer.from_pretrained('VidhuMathur/bart-log-summarization')

model = pipeline("summarization", model=model, tokenizer=tokenizer)

text = "summarize: 2023-11-15T19:39:02.238394189Z stderr F 2023-11-15 19:39:02,237 INFO [__main__] [server.py:32] [trace_id=6011fa67839c66d0d44542ec0f996416 span_id=8aed01d1fe2a3174 resource.service.name=00688f8f-1904-429a-80b9-06b2c92df17d trace_sampled=True] - executed query: SELECT * FROM profiles WHERE id = '1529' , time taken: 0:00:00.000541"
summary = model(text)

print(summary[0]['summary_text'])
