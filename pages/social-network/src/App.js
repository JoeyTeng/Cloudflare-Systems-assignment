import { Router } from "@reach/router";

import Posts from './components/posts'
import Publish from './components/publish'

function App() {
  return (
    <Router>
      <Posts path="/" />
      <Publish path="/publish" />
    </Router>
  );
}

export default App;

