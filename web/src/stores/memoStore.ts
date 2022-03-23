import utils from "../helpers/utils";

export interface State {
  memos: Model.Memo[];
  tags: string[];
}

interface SetMemosAction {
  type: "SET_MEMOS";
  payload: {
    memos: Model.Memo[];
  };
}

interface SetTagsAction {
  type: "SET_TAGS";
  payload: {
    tags: string[];
  };
}

interface InsertMemoAction {
  type: "INSERT_MEMO";
  payload: {
    memo: Model.Memo;
  };
}

interface DeleteMemoByIdAction {
  type: "DELETE_MEMO_BY_ID";
  payload: {
    id: string;
  };
}

interface EditMemoByIdAction {
  type: "EDIT_MEMO";
  payload: Model.Memo;
}

export type Actions = SetMemosAction | SetTagsAction | InsertMemoAction | DeleteMemoByIdAction | EditMemoByIdAction;

export function reducer(state: State, action: Actions): State {
  switch (action.type) {
    case "SET_MEMOS": {
      const memos = utils.dedupeObjectWithId(
        action.payload.memos.sort((a, b) => utils.getTimeStampByDate(b.createdAt) - utils.getTimeStampByDate(a.createdAt))
      );

      return {
        ...state,
        memos: [...memos],
      };
    }
    case "SET_TAGS": {
      return {
        ...state,
        tags: action.payload.tags,
      };
    }
    case "INSERT_MEMO": {
      const memos = utils.dedupeObjectWithId(
        [action.payload.memo, ...state.memos].sort((a, b) => utils.getTimeStampByDate(b.createdAt) - utils.getTimeStampByDate(a.createdAt))
      );

      return {
        ...state,
        memos,
      };
    }
    case "DELETE_MEMO_BY_ID": {
      return {
        ...state,
        memos: [...state.memos].filter((memo) => memo.id !== action.payload.id),
      };
    }
    case "EDIT_MEMO": {
      const memos = state.memos.map((m) => {
        if (m.id === action.payload.id) {
          return {
            ...m,
            ...action.payload,
          };
        } else {
          return m;
        }
      });

      return {
        ...state,
        memos,
      };
    }
    default: {
      return state;
    }
  }
}

export const defaultState: State = {
  memos: [],
  tags: [],
};
