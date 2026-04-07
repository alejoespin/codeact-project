TAREA: {user_request}
CONTEXTO:
- {context}
RESULTADOS:
- {observations}
REGLAS:
- NUNCA respondas directamente con tu conocimiento interno.
- Ten en cuenta la DOCUMENTACIÓN de las funciones mencionadas para implementar correctamente.
- Solo GENERA el código parcial ya que el proceso es iterativo y se debe resolver paso a paso.
- Responde con código Go ejecutable dentro de ```go ... ```.
- No agregues comentarios al código.
- El código debe retornar EXCLUSIVAMENTE el resultado sin textos adicionales.
- Responde EXCLUSIVAMENTE con el código en go, no agregues calculos diferentes.
- El RESULTADO del código se recibirá como OBSERVACIÓN.
- Si el resultado no es suficiente para completar la tarea, genera NUEVO código basado en la OBSERVACIÓN.
- REPITE hasta tener certeza del resultado final.
- SOLAMENTE emite FINAL: cuando tengas la certeza del resultado final.