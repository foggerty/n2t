;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;;; Functions to map VM tokens to assembler instructions.
;;;; Take a sequence of tokens representing VM instructions, and
;;;; output a sequence of assembly instructions.
;;;;
;;;; A 'token' is a map containing (at least) :cmd, and possibly :arg1
;;;; & :arg2.

(ns n2t.mapper
  (:require [n2t.commands :refer :all]
            [clojure.string :as str]))

(defn matching-token [str tokens]
  "Returns the first matching token (key) from the tokens map,
  otherwise nil.  ToDo - token is a confusing name, given the
  context..."
  (some #(if (str/includes? str %) % nil)
        (keys tokens)))

(defn replace-vars [replacements instructions]
  "Replacements should be a map of token/replacement.  Any instance
  will be replaced in instructions.  Note that there should only ever
  be ONE token to replace per instruction."
  (map #(let [replacement (matching-token % replacements)]
          (if replacement
            (str/replace % replacement (get replacements replacement))
            %))
       instructions))

(def commands
  "A collection of VM commands, and the associated assembly instructions
  that should be emitted.  Note that there are placeholders for some
  commands.

  Some assumptions are made for all commands:

  All binary arithmetic operations expect that the two values are
  already in D & {{temp}}, where {{temp}} is R5-R12.

  All arithmetic commands will leave the result in D."
  {:inc-sp            ["// Inc sp"
                       "@SP"
                       "M=M+1"]

   :dec-sp            ["// Dec sp"
                       "@SP"
                       "M=M-1"]

   :stack-to-D        ["// Copy stack to d"
                       "@SP"
                       "A=M"
                       "D=M"]

   :stack-to-tmp      ["// Copy stack to @{{temp}"
                       "@SP"
                       "A=M"
                       "D=M"
                       "@{{temp}}"
                       "M=D"]

   :D-to-stack        ["// Copy D to stack"
                       "@SP"
                       "A=M"
                       "M=D"]

   :tmp-to-stack      ["// Copy R{{temp}} to stack"
                       "@{{temp}}"
                       "D=M"
                       "@SP"
                       "A=M"
                       "M=D"]

   :constant-to-stack ["// push constant {{constant}}"
                       "@{{constant}}"
                       "D=A"
                       "@SP"
                       "A=M"
                       "M=D"]

   :add               ["// add"
                       "@{{temp}}"
                       "D=D+M"]})

(defn arithmetic-dispatch [cmd]
  "Dispatch to the correct function based on arity."
  (if (contains? vm-cmds-unary cmd)
    :arithmetic-unary
    :arithmetic-binary))

(defn memory-dispatch [segment]
  "Dispatch to the correct segment handling function."
  (cond
    (= segment :argument) :memory-argument
    (= segment :local)    :memory-local
    (= segment :static)   :memory-static
    (= segment :constant) :memory-constant
    (= segment :this)     :memory-this
    (= segment :that)     :memory-that
    (= segment :pointer)  :memory-pointer
    (= segment :temp)     :memory-temp))

(defn token-dispatch [token]
  (let [cmd (:cmd token)
        arg1 (:arg1 token)]
    (cond (contains? vm-cmds-control cmd) :control
          (contains? vm-cmds-function cmd) :function
          (contains? vm-cmds-arithmetic cmd) (arithmetic-dispatch cmd)
          (contains? vm-cmds-memory cmd) (memory-dispatch arg1))))

(defmulti map-token token-dispatch)

(defmethod map-token :arithmetic-binary [token]
  (replace-vars {"{{temp}}" "R5"}
                (mapcat #(% commands) [:dec-sp
                                       :stack-to-tmp
                                       :dec-sp
                                       :stack-to-D
                                       (:cmd token)
                                       :D-to-stack
                                       :inc-sp])))

(defmethod map-token :arithmetic-unary [token]
  (mapcat #(% commands) [:dec-sp
                         :copy-to-D
                         (:cmd token)
                         :D-to-stack
                         :inc-sp]))

(defmethod map-token :control [token])

(defmethod map-token :function [token])

(defmethod map-token :memory-argument [token])

(defmethod map-token :memory-local [token])

(defmethod map-token :memory-static [token])

(defmethod map-token :memory-constant [token]
  (replace-vars  {"{{constant}}" (:arg2 token)}
                 (mapcat #(% commands) [:constant-to-stack
                                        :inc-sp])))

(defmethod map-token :memory-this [token])

(defmethod map-token :memory-that [token])

(defmethod map-token :memory-pointer [token])

(defmethod map-token :memory-temp [token])

(defn tokens-to-asm [tokens]
  "It all seems so simple when viewed from here..."
  (mapcat #(let [src (str "// VM INSTRUCTION: " (str/upper-case (:source %)))]
             (conj (map-token %) src))
          tokens))
