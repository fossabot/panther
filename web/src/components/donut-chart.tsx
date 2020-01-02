import React from 'react';
import { Box, useTheme, defaultTheme } from 'pouncejs';

interface DonutChartDatum {
  color: keyof typeof defaultTheme['colors'];
  value: number;
  label: string;
}

interface DonutChartProps {
  /** A function that should return the value showcased in the middle of the donut */
  renderLabel: (data: DonutChartDatum[], index: number) => React.ReactNode;

  /** The data for the chart */
  data: DonutChartDatum[];
}

const DonutChart: React.FC<DonutChartProps> = ({ data, renderLabel }) => {
  const theme = useTheme();
  const container = React.useRef<HTMLDivElement>(null);

  React.useEffect(() => {
    // We are not allowed to put async function directly in useEffect. Instead, we should define
    // our own async function and call it within useEffect
    (async () => {
      // load the pie chart
      const [echarts] = await Promise.all([
        import(/* webpackChunkName: "echarts" */ 'echarts/lib/echarts'),
        import(/* webpackChunkName: "echarts" */ 'echarts/lib/chart/pie'),
        import(/* webpackChunkName: "echarts" */ 'echarts/lib/component/legend'),
      ]);

      // initialize a chart in the given DOM element
      const donutChart = echarts.init(container.current);

      // map the data to the shape that echarts expects
      const eChartsData = data.map(({ value, label, color }) => ({
        value,
        name: label,
        itemStyle: { color: theme.colors[color] },
      }));

      // draw the pie chart
      donutChart.setOption({
        legend: {
          bottom: 0,
          data: eChartsData.map(d => d.name),
          textStyle: {
            fontSize: theme.fontSizes[2] as number,
            color: theme.colors.grey400,
            fontFamily: theme.fonts.primary,
          },
          icon: 'circle',
        },
        series: [
          {
            type: 'pie',
            center: ['50%', '35%'],
            radius: ['47.5%', '70%'],
            avoidLabelOverlap: false,
            animation: false,
            label: {
              normal: {
                show: false,
                position: 'center',
                fontSize: 44,
                fontWeight: 'bold',
                fontFamily: theme.fonts.primary,
              },
              emphasis: {
                show: true,
                formatter: ({ dataIndex }) => renderLabel(data, dataIndex),
                rich: {
                  small: {
                    fontSize: theme.fontSizes[2],
                    fontWeight: 'bold',
                  },
                },
              },
            },
            labelLine: {
              normal: {
                show: false,
              },
            },
            data: eChartsData,
          },
        ],
      });
    })();
  }, []);

  return <Box innerRef={container} width="100%" height="100%" />;
};

export default React.memo(DonutChart);
