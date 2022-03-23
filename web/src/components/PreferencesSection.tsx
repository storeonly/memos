import { useContext } from "react";
import appContext from "../stores/appContext";
import { globalStateService, memoService } from "../services";
import { parseHtmlToRawText } from "../helpers/marked";
import { formatMemoContent } from "./Memo";
import "../less/preferences-section.less";

interface Props {}

const PreferencesSection: React.FC<Props> = () => {
  const { globalState } = useContext(appContext);
  const { useTinyUndoHistoryCache, shouldHideImageUrl, shouldSplitMemoWord, shouldUseMarkdownParser } = globalState;

  const demoMemoContent = "👋 Hiya, welcome to memos!\n* ✨ **Open source project**;\n* 😋 What do you think;\n* 📑 Tell me something plz;";

  const handleOpenTinyUndoChanged = () => {
    globalStateService.setAppSetting({
      useTinyUndoHistoryCache: !useTinyUndoHistoryCache,
    });
  };

  const handleSplitWordsValueChanged = () => {
    globalStateService.setAppSetting({
      shouldSplitMemoWord: !shouldSplitMemoWord,
    });
  };

  const handleHideImageUrlValueChanged = () => {
    globalStateService.setAppSetting({
      shouldHideImageUrl: !shouldHideImageUrl,
    });
  };

  const handleUseMarkdownParserChanged = () => {
    globalStateService.setAppSetting({
      shouldUseMarkdownParser: !shouldUseMarkdownParser,
    });
  };

  const handleExportBtnClick = async () => {
    const formatedMemos = memoService.getState().memos.map((m) => {
      return {
        content: m.content,
        createdAt: m.createdAt,
      };
    });

    const jsonStr = JSON.stringify(formatedMemos);
    const element = document.createElement("a");
    element.setAttribute("href", "data:text/json;charset=utf-8," + encodeURIComponent(jsonStr));
    element.setAttribute("download", "data.json");
    element.style.display = "none";
    document.body.appendChild(element);
    element.click();
    document.body.removeChild(element);
  };

  const handleFormatMemosBtnClick = async () => {
    const memos = memoService.getState().memos;
    for (const m of memos) {
      memoService.updateMemo(m.id, parseHtmlToRawText(m.content));
    }
  };

  return (
    <>
      <div className="section-container preferences-section-container">
        <p className="title-text">Memo Display</p>
        <div
          className="demo-content-container memo-content-text"
          dangerouslySetInnerHTML={{ __html: formatMemoContent(demoMemoContent) }}
        ></div>
        <label className="form-label checkbox-form-label hidden" onClick={handleSplitWordsValueChanged}>
          <span className="normal-text">Auto-space in English and Chinese</span>
          <img className="icon-img" src={shouldSplitMemoWord ? "/icons/checkbox-active.svg" : "/icons/checkbox.svg"} />
        </label>
        <label className="form-label checkbox-form-label" onClick={handleUseMarkdownParserChanged}>
          <span className="normal-text">Partial markdown format parsing</span>
          <img className="icon-img" src={shouldUseMarkdownParser ? "/icons/checkbox-active.svg" : "/icons/checkbox.svg"} />
        </label>
        <label className="form-label checkbox-form-label" onClick={handleHideImageUrlValueChanged}>
          <span className="normal-text">Hide image url</span>
          <img className="icon-img" src={shouldHideImageUrl ? "/icons/checkbox-active.svg" : "/icons/checkbox.svg"} />
        </label>
      </div>
      <div className="section-container preferences-section-container">
        <p className="title-text">Editor Extensions</p>
        <label className="form-label checkbox-form-label" onClick={handleOpenTinyUndoChanged}>
          <span className="normal-text">
            Use{" "}
            <a target="_blank" href="https://github.com/boojack/tiny-undo" onClick={(e) => e.stopPropagation()} rel="noreferrer">
              tiny-undo
            </a>
          </span>
          <img className="icon-img" src={useTinyUndoHistoryCache ? "/icons/checkbox-active.svg" : "/icons/checkbox.svg"} />
        </label>
      </div>
      <div className="section-container">
        <p className="title-text">Others</p>
        <div className="w-full flex flex-row justify-start items-center">
          <button className="px-2 py-1 border rounded text-base hover:opacity-80" onClick={handleExportBtnClick}>
            Export data as JSON
          </button>
          <button className="btn format-btn hidden" onClick={handleFormatMemosBtnClick}>
            Format Data
          </button>
        </div>
      </div>
    </>
  );
};

export default PreferencesSection;
