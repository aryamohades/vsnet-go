import React from 'react';
import { hot } from 'react-hot-loader';
import { throttle } from 'throttle-debounce';
import prettyBytes from 'pretty-bytes';
import './app.css';

class App extends React.Component {
  constructor(props) {
    super(props);
    this.fetchStats = throttle(3000, this.fetchStats);

    this.state = {
      data: null,
      error: null,
      dots: ''
    };
  }

  componentDidMount() {
    this.fetchStats();
    this.interval = this.updateDots();
  }

  componentWillUnmount() {
    clearInterval(this.interval);
  }

  async fetchStats() {
    try {
      const response = await fetch('http://172.17.0.2:8080/stats');
      const data = await response.json();

      this.setState({
        data,
        error: null
      });
    } catch (e) {
      this.setState({
        data: null,
        error: e
      });
    } finally {
      this.fetchStats();
    }
  }

  updateDots() {
    return setInterval(() => {
      this.setState(prevState => {
        let dots = prevState.dots + '.';

        if (dots.length > 4) {
          dots = '';
        }

        return { dots };
      });
    }, 400);
  }

  render() {
    const { data, error, dots } = this.state;

    return (
      <div className="app">
        <h2>Stats</h2>
        {data && (
          <>
            <div className="stat-container">
              <div className="stat-label">Connections</div>
              <div>{data.connections}</div>
            </div>
            <div className="stat-container">
              <div className="stat-label">Num Goroutines</div>
              <div>{data.num_goroutines}</div>
            </div>
            <div className="stat-container">
              <div className="stat-label">Heap Alloc</div>
              <div>{prettyBytes(data.heap_alloc)}</div>
            </div>
            <div className="stat-container">
              <div className="stat-label">Total Alloc</div>
              <div>{prettyBytes(data.total_alloc)}</div>
            </div>
            <div className="stat-container">
              <div className="stat-label">Heap Sys</div>
              <div>{prettyBytes(data.heap_sys)}</div>
            </div>
            <div className="stat-container">
              <div className="stat-label">GC Cycles</div>
              <div>{data.num_gc}</div>
            </div>
            <div className="stat-container">
              <div className="stat-label">Next GC</div>
              <div>{prettyBytes(data.next_gc)}</div>
            </div>
            <div className="stat-container">
              <div className="stat-label">GCCPU Fraction</div>
              <div>{data.gccpu_fraction}</div>
            </div>
            <div className="stat-container">
              <div className="stat-label">Pause Total</div>
              <div>{`${(data.pause_total_ns / 1000000).toFixed(2)} ms`}</div>
            </div>
          </>
        )}
        {error && (
          <div className="no-data">
            Waiting for data to become available<span>{dots}</span>
          </div>
        )}
      </div>
    );
  }
}

export default hot(module)(App);
