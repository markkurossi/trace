;;;
;;; Copyright (c) 2023 Markku Rossi
;;;
;;; All rights reserved.
;;;

;; Trace handler for multi-party computation (MPC) evaluation.
(library (main)
  (export)
  (import (rnrs lists))

  ;; Handle log record.
  (define (handle-record r)
    (display r) (newline))
  )
