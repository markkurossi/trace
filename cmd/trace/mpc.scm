;;;
;;; Copyright (c) 2023 Markku Rossi
;;;
;;; All rights reserved.
;;;

;; Trace handler for multi-party computation (MPC) evaluation.
(library (main)
  (export)
  (import (rnrs lists) (vt100 cursor) (vt100 graphics) (go format))

  (define (print-at row msg)
    (cursor-move row 0)
    (display msg)
    (erase-line-tail))

  ;; Handle log record.
  (define (handle-record r)
    (let ((message (cdr (assoc "message" r)))
          (level (cdr (assoc "level" r)))
          (time (cdr (assoc "time" r)))
          (attrs (cdr (assoc "attrs" r))))
      (case message
        (("progress") (msg-progress attrs))
        (else
         (print-at 2 message)
         (print-at 3 level)
         (print-at 4 time)
         (print-at 5 attrs)
         (print-at 6 "")
         ))))

  (define (msg-progress attrs)
    (let ((done (cdr (assoc "done" attrs)))
          (total (cdr (assoc "total" attrs)))
          (gates (cdr (assoc "gates" attrs))))
      (print-at 1 (hblock cols (/ (integer->float done) total)))
      (msg-progress-gates gates)))

  (define last-xor 0)
  (define last-xnor 0)
  (define last-and 0)
  (define last-or 0)
  (define last-inv 0)

  (define (msg-progress-gates gates)
    (let* ((num-xor (cdr (assoc "xor" gates)))
           (num-xnor (cdr (assoc "xnor" gates)))
           (num-and (cdr (assoc "and" gates)))
           (num-or (cdr (assoc "or" gates)))
           (num-inv (cdr (assoc "inv" gates)))

           (new-xor  (- num-xor  last-xor))
           (new-xnor (- num-xnor last-xnor))
           (new-and  (- num-and  last-and))
           (new-or   (- num-or   last-or))
           (new-inv  (- num-inv  last-inv))
           (num (+ new-xor new-xnor new-and new-or new-inv 0.0))
           (w (- cols 5)))

      (set! last-xor num-xor)
      (set! last-xnor num-xnor)
      (set! last-and num-and)
      (set! last-or num-or)
      (set! last-inv num-inv)

      (print-at 2 (string-append "xor  " (hblock cols (/ new-xor num))))
      (print-at 3 (string-append "xnor " (hblock cols (/ new-xnor num))))
      (print-at 4 (string-append "and  " (hblock cols (/ new-and num))))
      (print-at 5 (string-append "or   " (hblock cols (/ new-or num))))
      (print-at 6 (string-append "inv  " (hblock cols (/ new-inv num))))
      ))

  (define (redraw)
    (display "display: ") (display rows) (display "\xd7;") (display cols)
    (newline)
    )

  (erase-screen)
  )
