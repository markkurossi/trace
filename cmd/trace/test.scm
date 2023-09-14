(library (main)
  (export)
  (import (rnrs lists) (vt100 graphics))

  (display (assoc 'b '((a . 1) (b . 2) (c . 3)))) (newline)

  (define (loop i limit)
    (if (< i limit)
        (begin
          (display (hblock 20 (/ i limit))) (newline)
          (loop (+ i 1) limit))))

  (loop 0.0 20.0)
  )
