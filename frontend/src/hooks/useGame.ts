import { useContext } from "react";
import { GameContext } from "../app/providers/GameContext";

export function useGame() {
  const context = useContext(GameContext);
  if (!context) {
    throw new Error("useGame must be executed inside a GameProvider");
  }
  return context;
}
