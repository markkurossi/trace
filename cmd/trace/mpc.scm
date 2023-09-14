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
          (total (cdr (assoc "total" attrs))))
      (print-at 1 (hblock cols (/ (integer->float done) total)))))

  (define (redraw)
    (display "display: ") (display rows) (display "\xd7;") (display cols)
    (newline)
    )

  (erase-screen)
  )
