;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;;; Functions to map VM tokens to assembler instructions.
;;;; Takes a sequence of tokens, and outputs a sequence of instructions.

(ns mapper
  (:require [clojure.string :as str]))

;; push constant 1
;; @1
;; D=A
;;

(defn arg-count-ok? [token]
  (cond ()))

(defn token-to-asm [token]
  "Given a token (map) containing :cmd, :arg1 and :arg2, outputs the
  relevant ASM instructions."
  )
