import { forwardRef, ReactNode, useCallback, useContext, useEffect, useImperativeHandle, useRef } from "react";
import TinyUndo from "tiny-undo";
import appContext from "../../stores/appContext";
import { storage } from "../../helpers/storage";
import useRefresh from "../../hooks/useRefresh";
import Only from "../common/OnlyWhen";
import "../../less/editor.less";

export interface EditorRefActions {
  element: HTMLTextAreaElement;
  focus: FunctionType;
  insertText: (text: string) => void;
  setContent: (text: string) => void;
  getContent: () => string;
}

interface EditorProps {
  className: string;
  initialContent: string;
  placeholder: string;
  showConfirmBtn: boolean;
  showCancelBtn: boolean;
  tools?: ReactNode;
  onConfirmBtnClick: (content: string) => void;
  onCancelBtnClick: () => void;
  onContentChange: (content: string) => void;
}

// eslint-disable-next-line react/display-name
const Editor = forwardRef((props: EditorProps, ref: React.ForwardedRef<EditorRefActions>) => {
  const {
    globalState: { useTinyUndoHistoryCache },
  } = useContext(appContext);
  const {
    className,
    initialContent,
    placeholder,
    showConfirmBtn,
    showCancelBtn,
    onConfirmBtnClick: handleConfirmBtnClickCallback,
    onCancelBtnClick: handleCancelBtnClickCallback,
    onContentChange: handleContentChangeCallback,
  } = props;
  const editorRef = useRef<HTMLTextAreaElement>(null);
  const tinyUndoRef = useRef<TinyUndo | null>(null);
  // NOTE: auto-justify textarea height
  const refresh = useRefresh();

  useEffect(() => {
    if (!editorRef.current) {
      return;
    }

    if (initialContent) {
      editorRef.current.value = initialContent;
      refresh();
    }
  }, []);

  useEffect(() => {
    if (useTinyUndoHistoryCache) {
      if (!editorRef.current) {
        return;
      }

      const { tinyUndoActionsCache, tinyUndoIndexCache } = storage.get(["tinyUndoActionsCache", "tinyUndoIndexCache"]);

      tinyUndoRef.current = new TinyUndo(editorRef.current, {
        interval: 5000,
        initialActions: tinyUndoActionsCache,
        initialIndex: tinyUndoIndexCache,
      });

      tinyUndoRef.current.subscribe((actions, index) => {
        storage.set({
          tinyUndoActionsCache: actions,
          tinyUndoIndexCache: index,
        });
      });

      return () => {
        tinyUndoRef.current?.destroy();
      };
    } else {
      tinyUndoRef.current?.destroy();
      tinyUndoRef.current = null;
      storage.remove(["tinyUndoActionsCache", "tinyUndoIndexCache"]);
    }
  }, [useTinyUndoHistoryCache]);

  useEffect(() => {
    if (editorRef.current) {
      editorRef.current.style.height = "auto";
      editorRef.current.style.height = (editorRef.current.scrollHeight ?? 0) + "px";
    }
  }, [editorRef.current?.value]);

  useImperativeHandle(
    ref,
    () => ({
      element: editorRef.current as HTMLTextAreaElement,
      focus: () => {
        editorRef.current?.focus();
      },
      insertText: (rawText: string) => {
        if (!editorRef.current) {
          return;
        }

        const prevValue = editorRef.current.value;
        editorRef.current.value =
          prevValue.slice(0, editorRef.current.selectionStart) + rawText + prevValue.slice(editorRef.current.selectionStart);
        handleContentChangeCallback(editorRef.current.value);
        refresh();
      },
      setContent: (text: string) => {
        if (editorRef.current) {
          editorRef.current.value = text;
          handleContentChangeCallback(editorRef.current.value);
          refresh();
        }
      },
      getContent: (): string => {
        return editorRef.current?.value ?? "";
      },
    }),
    []
  );

  const handleEditorInput = useCallback(() => {
    handleContentChangeCallback(editorRef.current?.value ?? "");
    refresh();
  }, []);

  const handleEditorKeyDown = useCallback((event: React.KeyboardEvent<HTMLTextAreaElement>) => {
    event.stopPropagation();

    if (event.code === "Enter") {
      if (event.metaKey || event.ctrlKey) {
        handleCommonConfirmBtnClick();
      }
    }
    refresh();
  }, []);

  const handleCommonConfirmBtnClick = useCallback(() => {
    if (!editorRef.current) {
      return;
    }

    handleConfirmBtnClickCallback(editorRef.current.value);
    editorRef.current.value = "";
    refresh();
    // After confirm btn clicked, tiny-undo should reset state(clear actions and index)
    tinyUndoRef.current?.resetState();
  }, []);

  const handleCommonCancelBtnClick = useCallback(() => {
    handleCancelBtnClickCallback();
  }, []);

  return (
    <div className={"common-editor-wrapper " + className}>
      <textarea
        className="common-editor-inputer"
        rows={1}
        placeholder={placeholder}
        ref={editorRef}
        onInput={handleEditorInput}
        onKeyDown={handleEditorKeyDown}
      ></textarea>
      <div className="common-tools-wrapper">
        <div className="common-tools-container">
          <Only when={props.tools !== undefined}>{props.tools}</Only>
        </div>
        <div className="btns-container">
          <Only when={showCancelBtn}>
            <button className="action-btn cancel-btn" onClick={handleCommonCancelBtnClick}>
              Cancel editting
            </button>
          </Only>
          <Only when={showConfirmBtn}>
            <button className="action-btn confirm-btn" disabled={editorRef.current?.value === ""} onClick={handleCommonConfirmBtnClick}>
              Save <span className="icon-text">✍️</span>
            </button>
          </Only>
        </div>
      </div>
    </div>
  );
});

export default Editor;
