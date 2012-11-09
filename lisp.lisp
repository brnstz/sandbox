;; - What is Lisp?
;;   - Invented in 1958 by John McCarthy
;;   - Family of languages including Scheme, Clojure, Common Lisp
;;   - Influenced Javascript, Perl, Python, Ruby, many others
;;   - Functional programming basis, but also supports OOP and other paradigms
;;   - Contrast with Fortran, which influenced C, C++ and its descendants 
;;     (imperative/procedural programming)

;; - What is Common Lisp?
;;   - ANSI spec from 80s/90s
;;   - Unified different Lisp dialects 
;;   - 1980s AI boom
;;   - Compiled or interpreted / REPL
;;   - Garbage collection

;; - Tools
;;   - SBCL (Steel Bank Common Lisp)
;;   - Emacs
;;   - Slime (Superior Lisp Interation Mode for Emacs)

;; - Quotes 
;;  "Lisp is worth learning for the profound enlightenment experience you 
;;   will have when you finally get it. That experience will make you a better
;;   programmer for the rest of your days, even if you never actually use 
;;   Lisp itself a lot."
;;     - Eric Raymond, "How to Become a Hacker"

;;  "We were not out to win over the Lisp programmers; we were after the 
;;   C++ programmers. We managed to drag a lot of them about halfway to Lisp."
;;     - Guy Steele, Java spec co-author


;; - Functions, not operators

(+ 1 2 3 4 5 6)
(+ (- 6 5) (* 4 3) (/ 2 1))

(defvar plus-func #'+)
(defvar nums '(1 2 3 4 5 6))

(apply plus-func nums)

(eq plus-func #'+)

;; Define a function

(defun add-one (x)
  (let ((result (+ 1 x)))
    result))

;; Lambda

(lambda (x) (+ 1 x))

(mapcar #'(lambda (x) (+ 1 x)) '(1 2 3 4 5 6))

;; Dynamic but strong typing

;; wrong
(let ((x 1)
      (y "2"))
      
  (+ x y))

;; correct
(let ((x 1)
      (y 2))

  (+ x y))

;; Some data types
;;
;; - Symbols
;; - Lists
;; - Arbitrary size integers
;; - Hash tables
;; - Functions
;; - Others (arrays, packages, objects, etc.)

;; Recursion

(defun fib (i) 
  (cond
    ((= i 0) 0)
    ((< i 3) 1)
    (t (+ (fib (- i 1)) (fib (- i 2))))))

(fib 40)

;; Code is data

'(defun add-two (x)
  (+ 2 x))

(eval `(defun add-two (x) 
	 (+ 2 x)))

;; Macros

(defun print-trace (name &rest args)
  (format t "Calling ~a with: " name)
  (dolist (arg args) 
    (format t "~a " arg))
  (format t "~%"))

(defmacro defun-trace (name args body)
  `(defun ,name ,args 
     (print-trace ',name ,@args)
     ,body))

(defun-trace fibt (i) 
  (cond
    ((= i 0) 0)
    ((< i 3) 1)
    (t (+ (fibt (- i 1)) (fibt (- i 2))))))

;; Macroexpand

;; Memoization example

(defmacro defun-memo (name args body)
  (let ((cache (gensym))
	(retval (gensym))
	(argskey (cons 'list args)))
    `(let ((,cache (make-hash-table :test #'equalp)))
       (defun ,name ,args 
	 (or (gethash ,argskey ,cache)
	     (let ((,retval ,body))
	       (print-trace ',name ,@args)
	       (setf (gethash ,argskey ,cache) ,retval)
	       ,retval))))))
		 
(defun-memo fibm (i) 
  (cond
    ((= i 0) 0)
    ((< i 3) 1)
    (t (+ (fibm (- i 1)) (fibm (- i 2))))))

;; Further Reading
;;
;; - ANSI Common Lisp
;; - Successful Lisp 
;; - Practical Common Lisp
;; - Common Lisp Hyperspec
;; - comp.lang.lisp
;; - http://cliki.net (Common Lisp Wiki)
