import { createContext } from "svelte";
import type { Map as MapMap } from "maplibre-gl";

export const [getMapHandleContext, setMapHandleContext] = createContext<{
  instance?: MapMap;
}>();
