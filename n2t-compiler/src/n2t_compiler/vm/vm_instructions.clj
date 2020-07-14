;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;;; Maps VM string tokens to keywords.

(def vm-tokems
  "Maps VM string tokens to keywords."
  {;; Arithmetic and logic
   "add"       :add
   "sub"       :sub
   "neg"       :neg
   "eq"        :eq
   "gt"        :gt
   "lt"        :lt
   "and"       :and
   "or"        :or
   "not"       :not

   ;; Memory
   "push"      :push
   "pop"       :pop

   ;; Control flow
   "label"     :label
   "goto"      :goto
   "if-goto"   :if-goto

   ;; Function calling
   "function"  :function
   "call"      :call
   "return"    :return})

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;;; Organise into collections to make parsing easier

(def vm-arithmetic-cmds
  #{:add :sub :neg :eq :gt :lt :and :or :not})

(def vm-memory-cmds
  #{:push :pop})

(def vm-control-cmds
  #{:label :goto :if-goto})

(def vm-function-cmds
  #{:function :call :return})
